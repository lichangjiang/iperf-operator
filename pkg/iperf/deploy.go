package iperf

import (
	"fmt"
	"strconv"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	iperfalpha1clientset "github.com/lichangjiang/iperf-operator/pkg/client/clientset/versioned"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/lichangjiang/k8s/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

const (
	ServicePrefix = "service"
	ServerPrefix  = "iperf-server"
	ClientPrefix  = "iperf-client"
	RoleLabel     = "role"
	AppLabel      = "app"
)

type IperfTaskInfo struct {
	Namespace    string
	Name         string
	Uid          string
	ToEmail      string
	Image        string
	iperfTask    *iperfalpha1.IperfTask
	Port         int32
	ClientConfig ClientConfig
}

type ClientConfig struct {
	Interval int32
	Duration int32
}

type IperfTaskDeployer struct {
	k8sClient     kubernetes.Interface
	iperfClient   iperfalpha1clientset.Interface
	iperfTaskInfo *IperfTaskInfo
	ownerRef      metav1.OwnerReference
}

func NewIperfTaskInfo(iperfTask *iperfalpha1.IperfTask, namespace, name, uid string) *IperfTaskInfo {
	port := iperfTask.Spec.ServerSpec.Port
	if port == 0 {
		port = 5201
	}

	interval := iperfTask.Spec.ClientSpec.Interval
	if interval == 0 {
		interval = 2
	}

	duration := iperfTask.Spec.ClientSpec.Duration
	if duration == 0 {
		duration = 10
	}

	image := iperfTask.Spec.IperfImage
	email := iperfTask.Spec.ToEmail

	return &IperfTaskInfo{
		iperfTask: iperfTask,
		Image:     image,
		Port:      port,
		ToEmail:   email,
		Name:      name,
		Namespace: namespace,
		Uid:       uid,
		ClientConfig: ClientConfig{
			Interval: interval,
			Duration: duration,
		},
	}
}

func (info *IperfTaskInfo) GetIperfTask() *iperfalpha1.IperfTask {
	return info.iperfTask
}

func NewIperfTaskDeployer(k8sClient kubernetes.Interface,
	iperfClient iperfalpha1clientset.Interface,
	info *IperfTaskInfo) *IperfTaskDeployer {
	ownerRef := util.IperfTaskOwnRef(info.Namespace, info.Uid)
	return &IperfTaskDeployer{
		k8sClient:     k8sClient,
		iperfClient:   iperfClient,
		iperfTaskInfo: info,
		ownerRef:      ownerRef,
	}
}

func (deployer *IperfTaskDeployer) Run() (string, error) {
	nodesMap, err := kubeutil.GetNodeHostNames(deployer.k8sClient)
	if err != nil {
		return "", err
	}

	//每个node一个iperf server pod并且每个对应一个service
	serverIpMap, err := deployer.waitToCreateDeployAndSVC(nodesMap)
	if err != nil {
		return "", err
	}
	serverKeymap, statisMap := deployer.dispatchJobs(nodesMap, serverIpMap)
	if len(statisMap) == 0 {
		return "", fmt.Errorf("get empty iperf statis map")
	} else {
		content := HtmlTablePrint(serverKeymap, statisMap)
		klog.Infof("email content: %s", content)
		return content, nil
	}
}

func (deployer *IperfTaskDeployer) createServiceForDeployment(deployment *appsv1.Deployment) (string, error) {
	/*    deployName := deployment.ObjectMeta.Name*/
	//serviceName := ServerPrefix + "-" + deployName
	//labelValue := deployment.ObjectMeta.Labels[AppLabel]
	//labels := map[string]string{AppLabel: labelValue}
	//s := &corev1.Service{
	//ObjectMeta: metav1.ObjectMeta{
	//Name:   serviceName,
	//Labels: labels,
	//},
	//Spec: corev1.ServiceSpec{
	//Ports: []corev1.ServicePort{
	//{
	//Name:       serviceName,
	//Port:       deployer.iperfTaskInfo.Port,
	//TargetPort: intstr.FromInt(int(deployer.iperfTaskInfo.Port)),
	//},
	//},
	//},
	//}
	//kubeutil.SetOwnerRef(&s.ObjectMeta, &deployer.ownerRef)
	//s, err := deployer.k8sClient.CoreV1().Services(deployer.iperfTaskInfo.Namespace).Create(s)
	//if err != nil {
	//if !errors.IsAlreadyExists(err) {
	//return "", fmt.Errorf("failed to create server deployment service.%+v", err)
	//}
	//s, err = deployer.k8sClient.CoreV1().Services(deployer.iperfTaskInfo.Namespace).Get(serviceName, metav1.GetOptions{})
	//if err != nil {
	//return "", fmt.Errorf("failed to get server deployment service %s.%+v", serviceName, err)
	//}
	//}

	//if s == nil {
	//klog.Warningf("server deployment service %s ip not found.", serviceName)
	//return "", nil
	//}
	/*return s.Spec.ClusterIP, nil*/

	labelV := deployment.Labels[AppLabel]
	if labelV == "" {
		return "", fmt.Errorf("empty app label of deployment %s/%s", deployment.Namespace, deployment.Name)
	}

	return kubeutil.GetPodIpWithLabel(deployer.k8sClient, deployment.Namespace, AppLabel+"="+labelV, 60)
}

func (deployer *IperfTaskDeployer) createNodeDeployment(hostName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServerPrefix + "-" + hostName,
			Namespace: deployer.iperfTaskInfo.Namespace,
			Labels: map[string]string{
				AppLabel:  ServerPrefix + "-" + hostName,
				RoleLabel: ServerPrefix,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: util.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					AppLabel:  ServerPrefix + "-" + hostName,
					RoleLabel: ServerPrefix,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						AppLabel:  ServerPrefix + "-" + hostName,
						RoleLabel: ServerPrefix,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    ServerPrefix,
							Command: []string{"iperf3"},
							Args:    []string{"-s", "-p", strconv.Itoa(int(deployer.iperfTaskInfo.Port))},
							Image:   deployer.iperfTaskInfo.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: deployer.iperfTaskInfo.Port,
								},
							},
						},
					},
					NodeSelector: map[string]string{kubeutil.LabelHostname: hostName},
				},
			},
		},
	}
}

func (deployer *IperfTaskDeployer) waitToCreateDeployAndSVC(nodesMap map[string]string) (map[string]string, error) {
	serverIpMap := make(map[string]string)
	for _, v := range nodesMap {
		deployment := deployer.createNodeDeployment(v)
		kubeutil.SetOwnerRef(&deployment.ObjectMeta, &deployer.ownerRef)
		deployment, err := deployer.k8sClient.AppsV1().Deployments(deployer.iperfTaskInfo.Namespace).Create(deployment)
		if err != nil {
			return nil, fmt.Errorf("failed create deployment to node:%s. %+v", v, err)
		}
		ip, err := deployer.createServiceForDeployment(deployment)
		if err != nil {
			return nil, err
		}
		klog.Infof("server deployment for node %s ip is %s", v, ip)
		serverIpMap[v] = ip
	}

	err := kubeutil.WaitingForLabeledPodsToRun(deployer.k8sClient, RoleLabel+"="+ServerPrefix, deployer.iperfTaskInfo.Namespace, 300)
	if err != nil {
		return nil, err
	}

	return serverIpMap, nil
}

//顺序执行点对点iperf测试
func (deployer *IperfTaskDeployer) dispatchJobs(nodesMap map[string]string,
	serverIpMap map[string]string) (map[string][]CSKey,
	map[CSKey]IperfClientStatis) {

	serverKeyMap := make(map[string][]CSKey)
	statisMap := make(map[CSKey]IperfClientStatis)
	for _, serverv := range nodesMap {
		ip := serverIpMap[serverv]

		var csKeys []CSKey
		for _, v := range nodesMap {
			if v == serverv {
				continue
			}

			job := NewIperfJob(deployer.iperfTaskInfo.Namespace,
				v, deployer.iperfTaskInfo.Image,
				ip, deployer.iperfTaskInfo.Port,
				deployer.iperfTaskInfo.ClientConfig,
				deployer.ownerRef)

			log, err := job.Run(deployer.k8sClient)
			if err == nil {
				iperfJson, err := ParseLog(log)
				if err != nil {
					klog.Warningf("IperfJob for node %s parse log error.%+v", v, err)
				} else {
					key := CSKey{
						Server: serverv,
						Client: v,
					}
					statisMap[key] = iperfJson.Analyse()
					csKeys = append(csKeys, key)
				}
			} else {
				klog.Warningf("IperfJob for node %s error.%+v", v, err)
			}
		}
		serverKeyMap[serverv] = csKeys
	}
	return serverKeyMap, statisMap
}
