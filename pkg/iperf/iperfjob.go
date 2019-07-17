package iperf

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lichangjiang/k8s/kubeutil"
	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

const (
	JobPrefix = "job-"
	JobLable  = "job"
)

type IperfJob struct {
	Namespace    string
	Name         string
	JobNode      string
	Image        string
	ServerIp     string
	ServerPort   int32
	ClientConfig ClientConfig
	ownerRef     metav1.OwnerReference
}

func NewIperfJob(namespace, jobNode, image, serverIp string,
	serverPort int32,
	config ClientConfig,
	ownerRef metav1.OwnerReference) *IperfJob {

	name := JobPrefix + jobNode + "-" + serverIp
	return &IperfJob{
		Namespace:    namespace,
		Name:         name,
		JobNode:      jobNode,
		Image:        image,
		ServerIp:     serverIp,
		ServerPort:   serverPort,
		ClientConfig: config,
		ownerRef:     ownerRef,
	}

}

func (job *IperfJob) Run(k8sClient kubernetes.Interface) (string, error) {
	j := job.createJob()
	if err := kubeutil.RunReplaceableJob(k8sClient, j); err != nil {
		return "", fmt.Errorf("failed to start iperf-client job. %+v", err)
	}

	if err := kubeutil.WaitingForLabeledPodsToRun(k8sClient, JobLable+"="+job.Name, job.Namespace, 300); err != nil {
		klog.Warning(err.Error())
	}

	timeout := time.Duration(2*job.ClientConfig.Duration) * time.Second
	if err := kubeutil.WaitForJobCompletion(k8sClient, j, timeout); err != nil {
		return "", fmt.Errorf("failed to complete iperf-client job. %+v", err)
	}

	log, err := kubeutil.GetPodLog(k8sClient, job.Namespace, JobLable+"="+job.Name)
	if err != nil {
		return "", fmt.Errorf("failed to get iperf-client job log. %+v", err)
	}

	kubeutil.DeleteBatchJob(k8sClient, job.Namespace, job.Name, false)
	klog.Infof("iperf client job %s log \n%s", job.Name, log)
	return log, nil
}

func (job *IperfJob) createJob() *batch.Job {
	j := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.Name,
			Namespace: job.Namespace,
		},
		Spec: batch.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						JobLable: job.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Command: []string{"iperf3"},
							Args: []string{
								"-c",
								job.ServerIp,
								"-p",
								strconv.Itoa(int(job.ServerPort)),
								"-i",
								strconv.Itoa(int(job.ClientConfig.Interval)),
								"-t",
								strconv.Itoa(int(job.ClientConfig.Duration)),
								"-J",
							},
							Name:  "iper-client",
							Image: job.Image,
						},
					},
					NodeSelector:  map[string]string{kubeutil.LabelHostname: job.JobNode},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	kubeutil.SetOwnerRef(&j.ObjectMeta, &job.ownerRef)
	return j
}
