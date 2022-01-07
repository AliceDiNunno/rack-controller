package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) GetDebugNodes() ([]corev1.Node, *e.Error) {
	pods, err := k8s.Client.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return pods.Items, nil
}

func (k8s kubernetesInstance) GetNodes() ([]clusterDomain.Node, *e.Error) {
	nodes, err := k8s.Client.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return nodesToDomain(nodes.Items), nil
}

func (k8s kubernetesInstance) getNode(name string) (*corev1.Node, *e.Error) {
	node, err := k8s.Client.CoreV1().Nodes().Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrNodeNotFound)
	}

	return node, nil
}

func (k8s kubernetesInstance) GetNode(name string) (*clusterDomain.Node, *e.Error) {
	foundNode, err := k8s.getNode(name)

	if err != nil {
		return nil, err
	}

	node := nodeToDomain(foundNode)

	return node, nil
}

func nodesToDomain(nodes []corev1.Node) []clusterDomain.Node {
	var nodeList []clusterDomain.Node

	for _, node := range nodes {
		domainNode := nodeToDomain(&node)
		if domainNode != nil {
			nodeList = append(nodeList, *domainNode)
		}
	}

	return nodeList
}

func nodeToDomain(node *corev1.Node) *clusterDomain.Node {
	if node == nil {
		return nil
	}

	var taints = []clusterDomain.NodeTaint{}

	for _, taint := range node.Spec.Taints {
		taints = append(taints, clusterDomain.NodeTaint{
			Key:    taint.Key,
			Effect: string(taint.Effect),
			Since:  taint.TimeAdded.Time,
		})
	}

	condition := clusterDomain.NodeCondition{
		NetworkUnavailable: false,
		DiskPressure:       false,
		PidPressure:        false,
		Ready:              false,
		Taints:             taints,
	}

	for _, currentCondition := range node.Status.Conditions {
		if currentCondition.Type == corev1.NodeReady {
			condition.Ready = currentCondition.Status == corev1.ConditionTrue
		}
		if currentCondition.Type == corev1.NodeDiskPressure {
			condition.DiskPressure = currentCondition.Status == corev1.ConditionTrue
		}
		if currentCondition.Type == corev1.NodePIDPressure {
			condition.PidPressure = currentCondition.Status == corev1.ConditionTrue
		}
		if currentCondition.Type == corev1.NodeNetworkUnavailable {
			condition.NetworkUnavailable = currentCondition.Status == corev1.ConditionTrue
		}
	}

	return &clusterDomain.Node{
		Id:           string(node.UID),
		Name:         node.Name,
		CreationDate: node.CreationTimestamp.Time,
		Ip:           node.Status.Addresses[0].Address,
		Hardware: clusterDomain.NodeHardware{
			Cores:   node.Status.Capacity.Cpu().Value(),
			Storage: node.Status.Capacity.Storage().Value(),
			Memory:  node.Status.Capacity.Memory().Value(),
		},
		AvailableHardware: clusterDomain.NodeHardware{
			Cores:   node.Status.Allocatable.Cpu().Value(),
			Storage: node.Status.Allocatable.Storage().Value(),
			Memory:  node.Status.Allocatable.Memory().Value(),
		},
		OperatingSystem: clusterDomain.NodeOperatingSystem{
			OSType:         node.Status.NodeInfo.OperatingSystem,
			OSName:         node.Status.NodeInfo.OSImage,
			OSArchitecture: node.Status.NodeInfo.Architecture,
			KernelVersion:  node.Status.NodeInfo.KernelVersion,
		},
		Condition: condition,
	}
}
