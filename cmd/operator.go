package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	iperfalpha1clientset "github.com/lichangjiang/iperf-operator/pkg/client/clientset/versioned"
	iperfcontroller "github.com/lichangjiang/iperf-operator/pkg/controller"
	"github.com/lichangjiang/iperf-operator/pkg/util"
	"github.com/lichangjiang/kubecontroller"
	"github.com/lichangjiang/kubeutil"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

var OperatorCmd = &cobra.Command{
	Use:    "operator",
	Short:  "Runs the iperf operator for orchestrating and managing IperfTask in a Kubernetes cluster",
	Long:   `Runs the ipef operator for orchestrating and managing IperfTask in a Kubernetes cluster`,
	Hidden: true,
}

func init() {
	OperatorCmd.RunE = startOperator
}

func startOperator(cmd *cobra.Command, args []string) error {

	podNs := os.Getenv("POD_NAMESPACE")
	podName := os.Getenv("POD_NAME")
	runEnv := os.Getenv("IPERF_OPERATOR_RUNENV")
	emailUser := os.Getenv("IPERF_EMAIL_USER")
	emailPwd := os.Getenv("IPERF_EMAIL_PWD")
	emailSmtp := os.Getenv("IPERF_EMAIL_SMTP")
	emailPort := os.Getenv("IPERF_EMAIL_PORT")

	if podNs == "" || podName == "" {
		return fmt.Errorf("empty pod namespace or pod name")
	}

	if emailUser == "" || emailPwd == "" || emailSmtp == "" || emailPort == "" {
		return fmt.Errorf("email user password smtpServer or port empty")
	}

	iperfcontroller.PodName = podName
	iperfcontroller.PodNameSpace = podNs
	util.User = emailUser
	util.Smtp = emailSmtp
	util.Pwd = emailPwd
	p, err := strconv.Atoi(emailPort)
	if err != nil {
		return fmt.Errorf("convert email port string to int error,%+v", err)
	}
	util.Port = p

	var k8sClient kubernetes.Interface
	var iperfClient iperfalpha1clientset.Interface

	if runEnv != "" {
		k8sClient, iperfClient, err = util.GetAllClientsetsOut(kubeutil.DefaultConfigStr())
	} else {
		k8sClient, iperfClient, err = util.GetAllClientsetsInner()
	}

	if err != nil {
		return err
	}

	image, err := kubeutil.GetPodImage(k8sClient, podNs, podName, "iperf-operator")
	if err != nil {
		return fmt.Errorf("failed to get iperf-operator image wirh error %+v", err)
	}
	iperfcontroller.IperfOperatorImage = image

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	controller := iperfcontroller.NewIperfController(k8sClient, iperfClient)

	watcher := kubecontroller.NewWatcher(controller, iperfcontroller.IperfTaskResource, podNs, iperfClient.IperfAlpha1().RESTClient())
	watcher.Watch(&iperfalpha1.IperfTask{}, 1, stopChan)

	select {
	case <-signalChan:
		klog.Info("shutdown signal received,exiting")
		close(stopChan)
	}
	return nil
}
