package controller

import (
	"fmt"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	iperfalpha1clientset "github.com/lichangjiang/iperf-operator/pkg/client/clientset/versioned"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/lichangjiang/k8s/kubecontroller"
	"github.com/lichangjiang/k8s/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

const (
	Finished           string = "finished"
	Dispatched         string = "dispatched"
	Failed             string = "failed"
	serviceAccountName        = "iperf-operator"
)

var IperfOperatorImage string = "riverlcj/iperf:v0.3.1"
var PodNameSpace = ""
var PodName = ""

var IperfTaskResource = kubecontroller.CustomResource{
	Name:    "iperftask",
	Plural:  "iperftasks",
	Group:   util.Group,
	Version: util.Version,
	Scope:   apiextensionsv1beta1.NamespaceScoped,
	Kind:    util.Kind,
}

type IperfTaskController struct {
	k8sClient   kubernetes.Interface
	iperfClient iperfalpha1clientset.Interface
}

func NewIperfController(k8sClient kubernetes.Interface,
	iperfClient iperfalpha1clientset.Interface) *IperfTaskController {

	return &IperfTaskController{
		k8sClient:   k8sClient,
		iperfClient: iperfClient,
	}
}

func (c *IperfTaskController) AddFunc(obj interface{}) {
	iperfTask := obj.(*iperfalpha1.IperfTask)

	if iperfTask.Status.State == Finished {
		delIperfTask(iperfTask, c.iperfClient)
	} else if iperfTask.Status.State == "" {
		uuid := iperfTask.ObjectMeta.GetUID()
		name := iperfTask.ObjectMeta.GetName()
		ns := iperfTask.ObjectMeta.GetNamespace()

		klog.Infof("new iperfTask %s/%s,spec : %+v", ns, name, iperfTask.Spec)
		deployment := createIperfDeployment(ns, name, string(uuid))

		_, err := c.k8sClient.AppsV1().Deployments(ns).Create(deployment)
		if err != nil {
			klog.Errorf("failed to create deployment %s with error %s\n", deployment.ObjectMeta.Name, err.Error())
		} else {
			iperfTask.Status.State = Dispatched
			iperfTask.Status.Deploy = deployment.Name

			ns := iperfTask.ObjectMeta.GetNamespace()
			name := iperfTask.ObjectMeta.GetName()
			_, err := c.iperfClient.IperfAlpha1().IperfTasks(ns).Update(iperfTask)
			if err != nil {
				klog.Errorf("failed to update IperfTask %s with error %s\n", ns+"/"+name, err.Error())
			}
		}
	}
}

func (c *IperfTaskController) UpdateFunc(obj, newobj interface{}) {
	iperfTask := newobj.(*iperfalpha1.IperfTask)
	ns := iperfTask.ObjectMeta.GetNamespace()
	name := iperfTask.ObjectMeta.GetName()
	state := iperfTask.Status.State

	email := iperfTask.Spec.ToEmail
	msg := iperfTask.Status.Message
	klog.Infof("update iperfTask %s/%s status : %+v", ns, name, iperfTask.Status)
	if state == Finished {
		go func() {
			err := util.SendEmail(email, fmt.Sprintf("iperfTask %s/%s finished", ns, name), msg)
			if err != nil {
				klog.Infof("send email failed for finished IperfTask %s\n", ns+"/"+name)
			} else {
				klog.Infof("try to delete finished IperfTask %s\n", ns+"/"+name)
				delIperfTask(iperfTask, c.iperfClient)
			}
		}()
	} else if state == Failed {
		if email != "" {
			var content string
			if msg == "" {
				content = "<h1>Iperf Task Failed</h1>"
			} else {
				content = "<h1>" + msg + "</h1>"
			}
			go util.SendEmail(email, fmt.Sprintf("iperfTask %s/%s failed", ns, name), content)
		}
	}

}

func (c *IperfTaskController) DeleteFunc(obj interface{}) {
	iperfTask := obj.(*iperfalpha1.IperfTask)
	ns := iperfTask.ObjectMeta.GetNamespace()
	name := iperfTask.ObjectMeta.GetName()

	klog.Infof("finished IperfTask %s deleted\n", ns+"/"+name)
}

func delIperfTask(iperfTask *iperfalpha1.IperfTask, client iperfalpha1clientset.Interface) {
	ns := iperfTask.ObjectMeta.GetNamespace()
	name := iperfTask.ObjectMeta.GetName()
	err := client.IperfAlpha1().IperfTasks(ns).Delete(name, nil)
	if err != nil {
		klog.Errorf("fail to delete iperftask %s with error %s\n", ns+"/"+name, err.Error())
	}
}

func createIperfDeployment(ns, name, uuid string) *appsv1.Deployment {
	deploymentName := "iperf-operator-deploy-" + uuid

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: ns,
			Labels: map[string]string{
				"app": "iperf-operator-deploy",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: util.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "iperf-operator-deploy",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "iperf-operator-deploy",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{
						{
							Name:    "iperf-operator-deploy",
							Command: []string{"iperf-operator", "deploy"},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "IPERFTASK_NAME",
									Value: name,
								},
								corev1.EnvVar{
									Name:  "IPERFTASK_NAMESPACE",
									Value: ns,
								},
								corev1.EnvVar{
									Name:  "IPERFTASK_UID",
									Value: uuid,
								},
							},
							Image: IperfOperatorImage,
						},
					},
				},
			},
		},
	}

	ownerRef := util.IperfTaskOwnRef(ns, uuid)
	kubeutil.SetOwnerRef(&deployment.ObjectMeta, &ownerRef)
	return deployment
}
