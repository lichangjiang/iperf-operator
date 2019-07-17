package controller

import (
	"testing"
	"time"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/lichangjiang/k8s/kubeutil"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var IperfTestTask *iperfalpha1.IperfTask = &iperfalpha1.IperfTask{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test-task",
		Namespace: "default",
		Labels: map[string]string{
			"iperf": "test-task",
		},
	},
	Spec: iperfalpha1.IperfSpec{
		IperfImage: "172.17.8.101:30002/networkstatic/iperf3",
		ToEmail:    "305120108@qq.com",
		ServerSpec: iperfalpha1.IperfServerSpec{
			Port: 9000,
		},
		ClientSpec: iperfalpha1.IperfClientSpec{
			Interval: 2,
			Duration: 10,
		},
	},
}

func TestIperfDeploy(t *testing.T) {
	iperfTask := IperfTestTask
	k8sClient, myClient, err := util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	assert.NilError(t, err)

	iperfTask, err = myClient.IperfAlpha1().IperfTasks("default").Create(iperfTask)
	assert.NilError(t, err)

	//default serviceaccount  make deployment fail
	deployment := createIperfDeployment(iperfTask.Namespace, iperfTask.Name, string(iperfTask.GetUID()))
	deployment, err = k8sClient.AppsV1().Deployments("default").Create(deployment)
	assert.NilError(t, err)

	err = kubeutil.WaitingForLabeledPodsToRun(k8sClient, "app=iperf-operator-deploy", "default", 30)
	assert.NilError(t, err)

	deployment, err = k8sClient.AppsV1().Deployments("default").Get(deployment.Name, metav1.GetOptions{})
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(int(deployment.Status.AvailableReplicas), 0))

	err = myClient.IperfAlpha1().IperfTasks("default").Delete(iperfTask.Name, nil)
	assert.NilError(t, err)

	time.Sleep(1 * time.Second)

	_, err = k8sClient.AppsV1().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
	assert.Assert(t, is.Equal(errors.IsNotFound(err), true))
}
