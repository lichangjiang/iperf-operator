package util

import (
	"fmt"

	iperfalpha1clientset "github.com/lichangjiang/iperf-operator/pkg/client/clientset/versioned"
	"github.com/lichangjiang/k8s/kubeutil"
	"k8s.io/client-go/kubernetes"
)

func GetIperfClientset(outConfig string, inCluster bool) (iperfalpha1clientset.Interface, error) {
	config, err := kubeutil.GetK8sConfig(outConfig, inCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config. %+v", err)
	}

	clientset, err := iperfalpha1clientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create iperf client set. %+v", err)
	}

	return clientset, nil
}

func GetAllClientsetsInner() (kubernetes.Interface,
	iperfalpha1clientset.Interface, error) {
	k8sClient, err := kubeutil.CreateK8sClientset("", true)
	if err != nil {
		return nil, nil, err
	}

	iperfClient, err := GetIperfClientset("", true)
	if err != nil {
		return nil, nil, err
	}

	return k8sClient, iperfClient, nil
}

func GetAllClientsetsOut(outConfig string) (kubernetes.Interface,
	iperfalpha1clientset.Interface, error) {
	k8sClient, err := kubeutil.CreateK8sClientset(outConfig, false)
	if err != nil {
		return nil, nil, err
	}

	iperfClient, err := GetIperfClientset(outConfig, false)
	if err != nil {
		return nil, nil, err
	}

	return k8sClient, iperfClient, nil
}
