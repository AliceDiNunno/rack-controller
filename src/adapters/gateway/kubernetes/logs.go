package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	v1 "k8s.io/api/core/v1"
)

func (k8s kubernetesInstance) GetPodLogs(namespace string, podName string) (string, *e.Error) {
	data := k8s.Client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
		Timestamps: true,
	})

	str, err := data.Do(context.Background()).Raw()

	return string(str), e.Wrap(err)
}
