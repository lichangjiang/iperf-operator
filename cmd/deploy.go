package cmd

import (
	"fmt"
	"os"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	iperfalpha1clientset "github.com/lichangjiang/iperf-operator/pkg/client/clientset/versioned"
	"github.com/lichangjiang/iperf-operator/pkg/controller"
	"github.com/lichangjiang/iperf-operator/pkg/iperf"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var DeployCmd = &cobra.Command{
	Use:    "deploy",
	Short:  "deploy iperftask deployment to manage the task instance",
	Long:   "deploy iperftask deployment to manage the task instance",
	Hidden: true,
}

func init() {
	DeployCmd.RunE = deploy
}

func deploy(cmd *cobra.Command, args []string) error {
	ns := os.Getenv("IPERFTASK_NAMESPACE")
	name := os.Getenv("IPERFTASK_NAME")
	uid := os.Getenv("IPERFTASK_UID")
	if ns == "" || name == "" {
		return fmt.Errorf("empty iperftask namespace and name error")
	}

	k8sClient, iperfClient, err := util.GetAllClientsetsInner()
	if err != nil {
		return err
	}

	iperfTask, err := iperfClient.IperfAlpha1().IperfTasks(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	cuid := string(iperfTask.ObjectMeta.UID)
	if cuid != uid {
		updateIperfTaskState(iperfClient, iperfTask, controller.Failed, fmt.Sprintf("iperfTask instance uid error ,need %s but found %s", uid, cuid))
		return fmt.Errorf("iperfTask instance uid error ,need %s but found %s", uid, cuid)
	}

	email := iperfTask.Spec.ToEmail
	if email == "" {
		updateIperfTaskState(iperfClient, iperfTask, controller.Failed, "iperfTask email is empty")
		return fmt.Errorf("iperfTask email is empty")
	}

	iperfTaskInfo := iperf.NewIperfTaskInfo(iperfTask, ns, name, uid)
	iperfDeployer := iperf.NewIperfTaskDeployer(k8sClient, iperfClient, iperfTaskInfo)
	emailContent, err := iperfDeployer.Run()
	if err != nil {
		updateIperfTaskState(iperfClient, iperfTask, controller.Failed, err.Error())
		return err
	}
	updateIperfTaskState(iperfClient, iperfTask, controller.Finished, emailContent)
	return nil
}

func updateIperfTaskState(client iperfalpha1clientset.Interface, task *iperfalpha1.IperfTask, state, message string) {
	task.Status.State = state
	task.Status.Message = message
	_, err := client.IperfAlpha1().IperfTasks(task.Namespace).Update(task)
	if err != nil {
		klog.Warningf("update iperfTask %s/%s to failed state error.%+v", task.Namespace, task.Name, err)
	}
}
