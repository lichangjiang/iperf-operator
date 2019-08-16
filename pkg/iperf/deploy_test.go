package iperf

import (
	"strings"
	"testing"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/lichangjiang/k8s/kubeutil"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
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
		IperfImage: "networkstatic/iperf3",
		ToEmail:    "305120108@qq.com",
		ServerSpec: iperfalpha1.IperfServerSpec{
			Port: 9000,
		},
		ClientSpec: iperfalpha1.IperfClientSpec{
			Mode:     "fast",
			Parallel: 1,
			Interval: 1,
			Duration: 2,
		},
	},
}

func TestGetNode(t *testing.T) {
	k8sClient, _, err := util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	assert.NilError(t, err)

	var nodesMap map[string]string
	nodesMap, err = kubeutil.GetNodeHostNames(k8sClient)
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(len(nodesMap), 3))
	for _, v := range nodesMap {
		t.Logf("node host:%s", v)
	}
}

func TestIperfServerDeploy(t *testing.T) {
	klog.InitFlags(nil)
	iperfTask := IperfTestTask
	k8sClient, myClient, err := util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	assert.NilError(t, err)

	iperfTask, err = myClient.IperfAlpha1().IperfTasks("default").Create(iperfTask)
	assert.NilError(t, err)

	ici := NewIperfTaskInfo(iperfTask, iperfTask.Namespace, iperfTask.Name, string(iperfTask.UID))
	deployer := NewIperfTaskDeployer(k8sClient, myClient, ici)

	var nodesMap map[string]string
	nodesMap, err = kubeutil.GetNodeHostNames(k8sClient)
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(len(nodesMap), 7))

	serverIpMap, err := deployer.waitToCreateDeployAndSVC(nodesMap)
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(len(serverIpMap), 7))

	err = myClient.IperfAlpha1().IperfTasks("default").Delete(iperfTask.Name, nil)
	assert.NilError(t, err)
}

func TestDispatchJob(t *testing.T) {
	klog.InitFlags(nil)
	iperfTask := IperfTestTask
	k8sClient, myClient, err := util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	assert.NilError(t, err)

	iperfTask, err = myClient.IperfAlpha1().IperfTasks("default").Create(iperfTask)
	assert.NilError(t, err)

	ici := NewIperfTaskInfo(iperfTask, iperfTask.Namespace, iperfTask.Name, string(iperfTask.UID))
	deployer := NewIperfTaskDeployer(k8sClient, myClient, ici)

	var nodesMap map[string]string
	nodesMap, err = kubeutil.GetNodeHostNames(k8sClient)
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(len(nodesMap), 8))

	serverIpMap, err := deployer.waitToCreateDeployAndSVC(nodesMap)
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(len(serverIpMap), 8))

	csKeyMap, statisMap := deployer.dispatchJobs(nodesMap, serverIpMap)
	assert.Assert(t, is.Equal(len(csKeyMap), 8))
	assert.Assert(t, is.Equal(len(statisMap), 56))

	err = myClient.IperfAlpha1().IperfTasks("default").Delete(iperfTask.Name, nil)
	assert.NilError(t, err)
}

func TestDeployRun(t *testing.T) {
	klog.InitFlags(nil)
	iperfTask := IperfTestTask
	k8sClient, myClient, err := util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	assert.NilError(t, err)

	iperfTask, err = myClient.IperfAlpha1().IperfTasks("default").Create(iperfTask)
	assert.NilError(t, err)

	ici := NewIperfTaskInfo(iperfTask, iperfTask.Namespace, iperfTask.Name, string(iperfTask.UID))
	deployer := NewIperfTaskDeployer(k8sClient, myClient, ici)

	content, err := deployer.Run()
	assert.NilError(t, err)
	assert.Assert(t, is.Equal(strings.HasPrefix(content, "<table"), true))
	err = myClient.IperfAlpha1().IperfTasks("default").Delete(iperfTask.Name, nil)
	assert.NilError(t, err)
}
