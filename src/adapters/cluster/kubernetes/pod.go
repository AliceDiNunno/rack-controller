package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/davecgh/go-spew/spew"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) GetPods(namespace string) ([]clusterDomain.Pod, *e.Error) {
	pods, err := k8s.Client.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return podsToDomain(pods.Items), nil
}

func (k8s kubernetesInstance) getPod(namespace string, name string) (*corev1.Pod, *e.Error) {
	pod, err := k8s.Client.CoreV1().Pods(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrPodNotFound)
	}

	if pod == nil {
		return nil, e.Wrap(err).Append(ErrPodNotFound)
	}

	return pod, nil
}

func (k8s kubernetesInstance) GetPod(namespace string, name string) (*clusterDomain.Pod, *e.Error) {
	foundPod, err := k8s.getPod(namespace, name)

	if err != nil {
		return nil, err
	}

	pod := podToDomain(foundPod)

	return pod, nil
}

func (k8s kubernetesInstance) GetDebugPods(namespace string) ([]corev1.Pod, *e.Error) {
	pods, err := k8s.Client.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return pods.Items, nil
}

func (k8s kubernetesInstance) GetPodsOfADeployment(namespace string, deployment string) ([]clusterDomain.Pod, *e.Error) {
	pods, err := k8s.Client.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{LabelSelector: "app = " + deployment})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return podsToDomain(pods.Items), nil
}

func (k8s kubernetesInstance) GetPodsOfANode(node string) ([]clusterDomain.Pod, *e.Error) {
	pods, err := k8s.Client.CoreV1().Pods("").List(context.Background(), v1.ListOptions{FieldSelector: "spec.nodeName=" + node})

	spew.Dump(pods, err)

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return podsToDomain(pods.Items), nil
}

func (k8s kubernetesInstance) DeletePod(namespace string, podName string) *e.Error {
	err := k8s.Client.CoreV1().Pods(namespace).Delete(context.Background(), podName, v1.DeleteOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToDeleteRessource)
	}

	return nil
}

func podsToDomain(pods []corev1.Pod) []clusterDomain.Pod {
	var podList []clusterDomain.Pod

	for _, pod := range pods {
		domainPod := podToDomain(&pod)
		if domainPod != nil {
			podList = append(podList, *domainPod)
		}
	}

	return podList
}

func podToDomain(pod *corev1.Pod) *clusterDomain.Pod {
	if pod == nil {
		return nil
	}

	var readyProbe *clusterDomain.PodProbe = nil

	/*	if pod.Spec.Containers[0].ReadinessProbe != nil {
		readyProbe = &domain.PodProbe{
			MountPath:   pod.Spec.Containers[0].ReadinessProbe.HTTPGet.MountPath,
			Scheme: string(pod.Spec.Containers[0].ReadinessProbe.HTTPGet.Scheme),
			Port:   pod.Spec.Containers[0].ReadinessProbe.HTTPGet.Port.StrVal,
		}
	}*/

	var status = pod.Status.Phase
	if pod.DeletionTimestamp != nil {
		status = "Terminating"
	}

	var podConditions = clusterDomain.PodCondition{
		Initialized:     false,
		Ready:           false,
		ContainersReady: false,
		PodScheduled:    false,
	}

	for _, condition := range pod.Status.Conditions {
		if condition.Status == "True" && condition.Type == "Initialized" {
			podConditions.Initialized = true
		}
		if condition.Status == "True" && condition.Type == "Ready" {
			podConditions.Ready = true
		}
		if condition.Status == "True" && condition.Type == "ContainersReady" {
			podConditions.ContainersReady = true
		}
		if condition.Status == "True" && condition.Type == "PodScheduled" {
			podConditions.PodScheduled = true
		}
	}

	var restartCount = -1
	if len(pod.Status.ContainerStatuses) > 0 {
		restartCount = int(pod.Status.ContainerStatuses[0].RestartCount)
	}

	return &clusterDomain.Pod{
		Id:           string(pod.UID),
		Name:         pod.Name,
		CreationDate: pod.CreationTimestamp.Time,
		Status:       string(status),
		Image: clusterDomain.PodImage{
			Name:  pod.Spec.Containers[0].Name,
			Image: pod.Spec.Containers[0].Image,
		},
		ReadyProbe:   readyProbe,
		NodeName:     pod.Spec.NodeName,
		QoS:          string(pod.Status.QOSClass),
		Condition:    podConditions,
		RestartCount: restartCount,
		InternalIp:   pod.Status.PodIP,
	}
}
