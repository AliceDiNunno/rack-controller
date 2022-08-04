package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) GetPersistentVolumeList(namespace string) ([]string, *e.Error) {
	persistentVolumes, err := k8s.Client.CoreV1().PersistentVolumes().List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	var persistentVolumeNames []string
	for _, persistentVolume := range persistentVolumes.Items {
		persistentVolumeNames = append(persistentVolumeNames, persistentVolume.Name)
	}

	return persistentVolumeNames, nil
}

func (k8s kubernetesInstance) GetPersistentVolumeClaimList(namespace string) ([]string, *e.Error) {
	persistentVolumeClaims, err := k8s.Client.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	var persistentVolumeClaimNames []string
	for _, persistentVolumeClaim := range persistentVolumeClaims.Items {
		persistentVolumeClaimNames = append(persistentVolumeClaimNames, persistentVolumeClaim.Name)
	}

	return persistentVolumeClaimNames, nil
}

func (k8s kubernetesInstance) GetPersistentVolume(namespace string, name string) (clusterDomain.PersistentVolume, *e.Error) {
	persistentVolume, err := k8s.Client.CoreV1().PersistentVolumes().Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return clusterDomain.PersistentVolume{}, e.Wrap(err).Append(ErrPersistentVolumeNotFound)
	}

	return clusterDomain.PersistentVolume{
		Name: persistentVolume.Name,
	}, nil
}

func (k8S kubernetesInstance) GetPersistentVolumeClaim(namespace string, name string) (clusterDomain.PersistentVolumeClaim, *e.Error) {
	persistentVolumeClaim, err := k8S.Client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return clusterDomain.PersistentVolumeClaim{}, e.Wrap(err).Append(ErrPersistentVolumeClaimNotFound)
	}

	return clusterDomain.PersistentVolumeClaim{
		Name: persistentVolumeClaim.Name,
	}, nil
}

func (k8s kubernetesInstance) CreatePersistentVolume(namespace string, persistentVolume clusterDomain.PersistentVolume) *e.Error {
	quantityFromDomain := resource.NewQuantity(persistentVolume.StorageSize, resource.BinarySI)
	creationType := v12.HostPathDirectoryOrCreate

	volumeFromDomain := v12.PersistentVolume{
		ObjectMeta: v1.ObjectMeta{
			Name: persistentVolume.Name,
			Labels: map[string]string{
				"type": "local",
			},
		},
		Spec: v12.PersistentVolumeSpec{
			Capacity: v12.ResourceList{
				v12.ResourceStorage: *quantityFromDomain,
			},
			PersistentVolumeSource: v12.PersistentVolumeSource{
				HostPath: &v12.HostPathVolumeSource{
					Path: persistentVolume.MountPath,
					Type: &creationType,
				},
			},
			AccessModes: []v12.PersistentVolumeAccessMode{
				v12.ReadWriteOnce,
			},
		},
	}

	_, err := k8s.Client.CoreV1().PersistentVolumes().Create(context.Background(), &volumeFromDomain, v1.CreateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToCreateRessource)
	}

	return nil
}

func (k8s kubernetesInstance) CreatePersistentVolumeClaim(namespace string, persistentVolumeClaim clusterDomain.PersistentVolumeClaim) *e.Error {
	quantityFromDomain := resource.NewQuantity(persistentVolumeClaim.StorageSize, resource.BinarySI)

	currentStorageClass := "nfs-client"
	//check if storage class contains nfs client
	_, err := k8s.Client.StorageV1().StorageClasses().Get(context.Background(), currentStorageClass, v1.GetOptions{})

	if err != nil {
		currentStorageClass = "local-path"
	}

	volumeClaimFromDomain := v12.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name: persistentVolumeClaim.Name,
		},
		Spec: v12.PersistentVolumeClaimSpec{
			AccessModes: []v12.PersistentVolumeAccessMode{
				v12.ReadWriteOnce,
			},
			Resources: v12.ResourceRequirements{
				Requests: v12.ResourceList{
					v12.ResourceStorage: *quantityFromDomain,
				},
			},
			StorageClassName: &currentStorageClass,
		},
	}

	_, err = k8s.Client.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), &volumeClaimFromDomain, v1.CreateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToCreateRessource)
	}

	return nil
}

/*
func (k8s kubernetesInstance) CreateService(namespace string, request clusterDomain.Service) *e.Error {
	service := v12.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      request.Name,
			Namespace: namespace,
		},
		Spec: v12.ServiceSpec{
			Type: v12.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": request.DeploymentName,
			},
			Ports: []v12.ServicePort{
				v12.ServicePort{
					Name:     request.PortName,
					Port:     request.Port,
					Protocol: v12.Protocol(request.Protocol),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(request.TargetPort),
					},
					NodePort: request.NodePort,
				},
			},
		},
	}

	_, err := k8s.Client.CoreV1().Services(namespace).Create(context.Background(), &service, v1.CreateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToCreateService)
	}

	return nil
}
*/
