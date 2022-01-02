package kubernetes

import (
	"context"
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/gateway/kubernetes/templates"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/davecgh/go-spew/spew"
	appsv1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1Apply "k8s.io/client-go/applyconfigurations/apps/v1"
	"time"
)

func (k8s kubernetesInstance) ListDeployments(namespace string) ([]domain.Deployment, *e.Error) {
	deployment, err := k8s.Client.AppsV1().Deployments(namespace).List(context.Background(), v12.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	return deploymentsToDomain(deployment.Items), nil
}

func (k8s kubernetesInstance) getDeployment(namespace string, name string) (*appsv1.Deployment, *e.Error) {
	deployment, err := k8s.Client.AppsV1().Deployments(namespace).Get(context.Background(), name, v12.GetOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrDeploymentNotFound)
	}

	return deployment, nil
}

func (k8s kubernetesInstance) GetDeployment(namespace string, name string) (*domain.Deployment, *e.Error) {
	foundDeployment, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return nil, err
	}

	deployment := deploymentToDomain(foundDeployment)

	return deployment, nil
}

func (k8s kubernetesInstance) GetDebugDeployments(namespace string) ([]appsv1.Deployment, *e.Error) {
	deployments, err := k8s.Client.AppsV1().Deployments(namespace).List(context.Background(), v12.ListOptions{})
	return deployments.Items, e.Wrap(err)
}

func (k8s kubernetesInstance) DeleteDeployment(namespace string, name string) *e.Error {
	_, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return err
	}

	deletionError := k8s.Client.AppsV1().Deployments(namespace).Delete(context.Background(), name, v12.DeleteOptions{})

	if err != nil {
		return e.Wrap(deletionError).Append(ErrUnableToDeleteRessource)
	}

	return nil
}

func (k8s kubernetesInstance) deployApp(namespace string, deployments []appsv1Apply.DeploymentApplyConfiguration) *e.Error {
	for _, deployment := range deployments {

		_, err := k8s.Client.AppsV1().Deployments(namespace).Apply(context.Background(), &deployment, v12.ApplyOptions{FieldManager: "rack-controller"})

		spew.Dump(err)

		if err != nil {
			return e.Wrap(ErrUnableToDeployApp)
		}
	}

	return nil
}

func (k8s kubernetesInstance) RestartDeployment(namespace string, name string) *e.Error {
	currentDeployment, err := k8s.Client.AppsV1().Deployments(namespace).Get(context.Background(), name, v12.GetOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrDeploymentNotFound)
	}

	if currentDeployment.Spec.Template.Annotations == nil {
		currentDeployment.Spec.Template.Annotations = map[string]string{}
	}
	currentDeployment.Spec.Template.Annotations["restarted"] = time.Now().String()

	_, err = k8s.Client.AppsV1().Deployments(namespace).Update(context.Background(), currentDeployment, v12.UpdateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToUpdateApp)
	}

	return nil
}

func (k8s kubernetesInstance) GetEnvironmentOfADeployment(namespace string, name string) ([]domain.Environment, *e.Error) {
	deployments, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return nil, err
	}

	containers := deployments.Spec.Template.Spec.Containers

	if containers == nil || len(containers) <= 0 {
		return nil, e.Wrap(ErrUnableToGetRessource)
	}

	env := containers[0].Env

	var envToReturn = []domain.Environment{}

	for _, envEntry := range env {
		newValue := domain.Environment{
			Name:  envEntry.Name,
			Value: envEntry.Value,
		}

		envToReturn = append(envToReturn, newValue)
	}

	return envToReturn, nil
}

func (k8s kubernetesInstance) GetPortsOfADeployment(namespace string, name string) ([]domain.Port, *e.Error) {
	deployments, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return nil, err
	}

	containers := deployments.Spec.Template.Spec.Containers

	if containers == nil || len(containers) <= 0 {
		return nil, e.Wrap(ErrUnableToGetRessource)
	}

	ports := containers[0].Ports

	var portsToReturn = []domain.Port{}

	service, err := k8s.getExposedPorts(namespace, fmt.Sprintf("%s", name))
	for _, portEntry := range ports {
		newValue := domain.Port{
			Name:            portEntry.Name,
			NetworkProtocol: string(portEntry.Protocol),
			ServicePort:     portEntry.ContainerPort,
		}

		if err == nil {
			for _, servicePort := range service.Spec.Ports {
				if servicePort.Port == portEntry.ContainerPort {
					newValue.ExposedPort = servicePort.NodePort
					break
				}
			}
		}

		portsToReturn = append(portsToReturn, newValue)
	}

	return portsToReturn, nil
}

func (k8s kubernetesInstance) GetConfigMapsOfADeployment(namespace string, name string) ([]string, *e.Error) {
	deployments, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return nil, err
	}

	containers := deployments.Spec.Template.Spec.Containers

	if containers == nil || len(containers) <= 0 {
		return nil, e.Wrap(ErrUnableToGetRessource)
	}

	configMaps := containers[0].EnvFrom

	var configMapsToReturn = []string{}

	for _, configMapEntry := range configMaps {
		if configMapEntry.ConfigMapRef != nil {
			configMapsToReturn = append(configMapsToReturn, configMapEntry.ConfigMapRef.Name)
		}
	}

	return configMapsToReturn, nil
}

func (k8s kubernetesInstance) GetSecretsOfADeployment(namespace string, name string) ([]string, *e.Error) {
	deployments, err := k8s.getDeployment(namespace, name)

	if err != nil {
		return nil, err
	}

	containers := deployments.Spec.Template.Spec.Containers

	if containers == nil || len(containers) <= 0 {
		return nil, e.Wrap(ErrUnableToGetRessource)
	}

	secrets := containers[0].EnvFrom

	var secretsToReturn = []string{}

	for _, secretsEntry := range secrets {
		if secretsEntry.SecretRef != nil {
			secretsToReturn = append(secretsToReturn, secretsEntry.SecretRef.Name)
		}
	}

	return secretsToReturn, nil
}

func (k8s kubernetesInstance) handleDeployment(namespace string, data interface{}) *e.Error {
	appDeployment, ok := data.(appsv1Apply.DeploymentApplyConfiguration)

	if ok {
		err := k8s.deployApp(namespace, []appsv1Apply.DeploymentApplyConfiguration{appDeployment})

		if err != nil {
			spew.Dump(err)
			return err
		}
	} else {
		return e.Wrap(ErrUnableToDeployApp)
	}

	return nil
}

func (k8s kubernetesInstance) CreateDeployment(namespace string, request domain.DeploymentCreationRequest) *e.Error {
	for _, template := range templates.AvailableTemplates() {
		if template.Id == request.TemplateId {
			template := template.Exec(namespace, request)

			err := k8s.handleDeployment(namespace, template)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deploymentsToDomain(deployments []appsv1.Deployment) []domain.Deployment {
	var deploymentList []domain.Deployment

	for _, deployment := range deployments {
		domainDeployment := deploymentToDomain(&deployment)
		if domainDeployment != nil {
			deploymentList = append(deploymentList, *domainDeployment)
		}
	}

	return deploymentList
}

func deploymentToDomain(deployment *appsv1.Deployment) *domain.Deployment {
	if deployment == nil {
		return nil
	}

	var replicas int64 = 1

	if deployment.Spec.Replicas != nil {
		replicas = int64(*deployment.Spec.Replicas)
	}

	pod := deployment.Spec.Template.Spec.Containers[0]

	var probe *domain.ContainerProbe

	/*
		if pod.ReadinessProbe != nil {
			probe = &domain.ContainerProbe{
				Path:   pod.ReadinessProbe.HTTPGet.Path,
				Scheme: string(pod.ReadinessProbe.HTTPGet.Scheme),
				Port:   pod.ReadinessProbe.HTTPGet.Port.StrVal,
			}
		}*/

	condition := domain.DeploymentCondition{}

	for _, currentCondition := range deployment.Status.Conditions {
		if currentCondition.Type == "Available" {
			condition.Available = currentCondition.Status == "True"
		}
		if currentCondition.Type == "Progressing" {
			condition.Progressing = currentCondition.Status == "True"
		}
	}

	return &domain.Deployment{
		Id:                string(deployment.UID),
		Name:              deployment.Name,
		ImageName:         deployment.Spec.Template.Spec.Containers[0].Image,
		CreationDate:      deployment.CreationTimestamp.Time,
		Generation:        deployment.Generation,
		Replicas:          replicas,
		UpdatedReplicas:   int64(deployment.Status.UpdatedReplicas),
		AvailableReplicas: int64(deployment.Status.AvailableReplicas),
		ReadyReplicas:     int64(deployment.Status.ReadyReplicas),
		Container: domain.Container{
			Name: pod.Name,
			Image: domain.ContainerImage{
				Name:  pod.Name,
				Image: pod.Image,
			},
			ReadyProbe: probe,
		},
		Condition: condition,
	}
}
