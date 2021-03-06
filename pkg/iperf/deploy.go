package iperf

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/lichangjiang/iperf-operator/pkg/algorithm"
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
	Mode     string
	Parallel int32
	Interval int32
	Duration int32
	Udp      bool
	BwLimit  string
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

	mode := iperfTask.Spec.ClientSpec.Mode
	if mode != "slow" {
		mode = "fast"
	} else {
		mode = "slow"
	}

	parallel := iperfTask.Spec.ClientSpec.Parallel
	if parallel < 1 {
		parallel = 1
	} else if parallel > 100 {
		parallel = 100
	}

	return &IperfTaskInfo{
		iperfTask: iperfTask,
		Image:     image,
		Port:      port,
		ToEmail:   email,
		Name:      name,
		Namespace: namespace,
		Uid:       uid,
		ClientConfig: ClientConfig{
			Mode:     mode,
			Parallel: parallel,
			Interval: interval,
			Duration: duration,
			Udp:      iperfTask.Spec.ClientSpec.Udp,
			BwLimit:  iperfTask.Spec.ClientSpec.BwLimit,
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
	nodesMap, err := kubeutil.GetNodeHostNamesWithFilter(deployer.k8sClient, func(node *corev1.Node) bool {
		if node.Spec.Unschedulable {
			return false
		}

		if _, hasMasterRoleLabel := node.Labels["node-role.kubernetes.io/master"]; hasMasterRoleLabel {
			return false
		}

		if len(node.Status.Conditions) == 0 {
			return false
		}

		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status != corev1.ConditionTrue {
				return false
			}
		}
		return true
	})
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
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errMsg := ""
	for _, v := range nodesMap {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			v := host
			deployment := deployer.createNodeDeployment(v)
			kubeutil.SetOwnerRef(&deployment.ObjectMeta, &deployer.ownerRef)
			deployment, err := deployer.k8sClient.AppsV1().Deployments(deployer.iperfTaskInfo.Namespace).Create(deployment)
			if err != nil {
				mutex.Lock()
				defer mutex.Unlock()
				errMsg = errMsg + fmt.Sprintf("failed create deployment to node:%s. %+v\n", v, err)
				return
			}
			ip, err := deployer.createServiceForDeployment(deployment)
			if err != nil {
				mutex.Lock()
				defer mutex.Unlock()
				errMsg = errMsg + err.Error() + "\n"
				return
			}
			klog.Infof("server deployment for node %s ip is %s", v, ip)
			mutex.Lock()
			serverIpMap[v] = ip
			mutex.Unlock()
		}(v)
	}
	wg.Wait()

	if errMsg != "" {
		return nil, fmt.Errorf(errMsg)
	}

	err := kubeutil.WaitingForLabeledPodsToRun(deployer.k8sClient, RoleLabel+"="+ServerPrefix, deployer.iperfTaskInfo.Namespace, 300)
	if err != nil {
		return nil, err
	}

	return serverIpMap, nil
}

type logMessage struct {
	server string
	client string
	statis IperfClientStatis
	err    error
}

//顺序执行点对点iperf测试
//nodesMap node名字 -> node HostName映射
//serverIpMap node HostName -> server ip映射
func (deployer *IperfTaskDeployer) dispatchJobs(nodesMap map[string]string,
	serverIpMap map[string]string) (map[string][]CSKey,
	map[CSKey]IperfClientStatis) {

	nodeLen := len(nodesMap)
	parallel := int(deployer.iperfTaskInfo.ClientConfig.Parallel)
	if parallel > nodeLen {
		parallel = nodeLen
	}

	mode := deployer.iperfTaskInfo.ClientConfig.Mode
	var jobMap algorithm.JobMap
	if parallel == 1 && mode == "fast" {
		jobMap = algorithm.FastSerialize(nodesMap, serverIpMap)
	} else if parallel > 1 {
		jobMap = algorithm.Parallelize(nodesMap, serverIpMap, parallel)
	} else {
		jobMap = algorithm.Serialize(nodesMap, serverIpMap)
	}

	serverKeyMap := make(map[string][]CSKey)
	statisMap := make(map[CSKey]IperfClientStatis)

	for i := 0; i < jobMap.EpochSize; i++ {
		jobs := jobMap.Jobs[i]
		jobLen := len(jobs)
		inChan := make(chan logMessage, 0)
		done := 0
		start := time.Now()
		for _, job := range jobs {
			go func(jobnode algorithm.JobNode) {
				logMsg := logMessage{
					server: jobnode.ServerHost,
					client: jobnode.ClientHost,
				}
				job := NewIperfJob(deployer.iperfTaskInfo.Namespace,
					jobnode.ClientHost, deployer.iperfTaskInfo.Image,
					jobnode.ServerIp, deployer.iperfTaskInfo.Port,
					deployer.iperfTaskInfo.ClientConfig,
					deployer.ownerRef)

				log, err := job.Run(deployer.k8sClient)
				if err == nil {
					iperfJson, err := ParseLog(log)
					if err != nil {
						logMsg.err = err
					} else {
						logMsg.statis = iperfJson.Analyse()
					}
				} else {
					logMsg.err = err
				}
				inChan <- logMsg
			}(job)
		}

		for done < jobLen {
			select {
			case msg := <-inChan:
				done++
				if msg.err == nil {
					key := CSKey{
						Server: msg.server,
						Client: msg.client,
					}
					statisMap[key] = msg.statis
					csKeys, ok := serverKeyMap[msg.server]
					if !ok {
						csKeys = []CSKey{}
					}
					serverKeyMap[msg.server] = append(csKeys, key)
				} else {
					klog.Infoln("----------------------------------------------")
					klog.Warningf("IperfJob for node %s error.%+v", msg.client, msg.err)
					klog.Infoln("----------------------------------------------")
				}
			}
		}
		elapsed := time.Now().Sub(start)
		klog.Infoln("**********************************************")
		klog.Infof("epoch %d all %d jobs finished in %f seconds\n", i, jobLen, elapsed.Seconds())
		klog.Infoln("**********************************************")
	}

	return serverKeyMap, statisMap
}
