package templates

import (
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes/utils"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/davecgh/go-spew/spew"
	corev1 "k8s.io/api/core/v1"
	appsv1Apply "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1Apply "k8s.io/client-go/applyconfigurations/core/v1"
	metav1Apply "k8s.io/client-go/applyconfigurations/meta/v1"
)

func ExecBasicContainerTemplate(namespace string, request clusterDomain.DeploymentCreationRequest) interface{} {
	apiVersion := "apps/v1"
	deploymentKind := "Deployment"
	appNamespace := namespace
	replicas := int32(request.Replicas)
	terminationGracePeriod := int64(30)
	pullPolicy := corev1.PullAlways
	ports := utils.PortsFromDomain(request.Ports)
	environment := utils.EnvironmentVariablesFromDomain(request.Environment)

	var envFrom []corev1Apply.EnvFromSourceApplyConfiguration

	spew.Dump(request.ConfigMaps)

	for _, currentConfigMap := range request.ConfigMaps {
		spew.Dump(currentConfigMap)
		name := currentConfigMap
		ref := corev1Apply.ConfigMapEnvSourceApplyConfiguration{
			LocalObjectReferenceApplyConfiguration: corev1Apply.LocalObjectReferenceApplyConfiguration{
				Name: &name,
			},
		}

		apply := corev1Apply.EnvFromSourceApplyConfiguration{
			Prefix:       nil,
			ConfigMapRef: &ref,
			SecretRef:    nil,
		}
		spew.Dump(apply)
		envFrom = append(envFrom, apply)
	}

	spew.Dump(envFrom)

	for _, currentSecret := range request.Secrets {
		apply := corev1Apply.EnvFromSourceApplyConfiguration{
			Prefix:       nil,
			ConfigMapRef: nil,
			SecretRef: &corev1Apply.SecretEnvSourceApplyConfiguration{
				LocalObjectReferenceApplyConfiguration: corev1Apply.LocalObjectReferenceApplyConfiguration{
					Name: &currentSecret,
				},
			},
		}
		envFrom = append(envFrom, apply)
	}

	createDeploymentConfiguration := appsv1Apply.DeploymentApplyConfiguration{
		TypeMetaApplyConfiguration: metav1Apply.TypeMetaApplyConfiguration{
			Kind:       &deploymentKind,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &metav1Apply.ObjectMetaApplyConfiguration{
			Name:      &request.DeploymentName,
			Namespace: &appNamespace,
		},
		Spec: &appsv1Apply.DeploymentSpecApplyConfiguration{
			Replicas: &replicas,
			Selector: &metav1Apply.LabelSelectorApplyConfiguration{
				MatchLabels:      map[string]string{"app": request.DeploymentName},
				MatchExpressions: nil,
			},
			Template: &corev1Apply.PodTemplateSpecApplyConfiguration{
				ObjectMetaApplyConfiguration: &metav1Apply.ObjectMetaApplyConfiguration{
					Labels: map[string]string{"app": request.DeploymentName},
				},
				Spec: &corev1Apply.PodSpecApplyConfiguration{
					Volumes: nil,
					Containers: []corev1Apply.ContainerApplyConfiguration{
						{
							Name:           &request.DeploymentName,
							Image:          &request.ImageName,
							Ports:          ports,
							Env:            environment,
							EnvFrom:        envFrom,
							VolumeMounts:   nil,
							VolumeDevices:  nil,
							LivenessProbe:  nil,
							ReadinessProbe: nil,
							Resources: &corev1Apply.ResourceRequirementsApplyConfiguration{
								/*							Limits: &corev1.ResourceList{
															"cpu": resource.Quantity{
																Format: resource.Format(fmt.Sprintf("%dM", request.CPU)),
															},
															"memory": resource.Quantity{
																Format: resource.Format(fmt.Sprintf("%dMi", request.Memory)),
															},
														},*/
								Requests: nil,
							},
							StartupProbe:    nil,
							ImagePullPolicy: &pullPolicy,
						},
					},
					TerminationGracePeriodSeconds: &terminationGracePeriod,
				},
			},
		},
		Status: nil,
	}

	return createDeploymentConfiguration
}
