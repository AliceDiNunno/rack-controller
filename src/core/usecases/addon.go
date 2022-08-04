package usecases

import (
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/google/uuid"

	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
)

func (i interactor) GetAddons(service *domain.Service) ([]domain.Addon, *e.Error) {
	addons, err := i.addonRepository.GetAddons(service)
	if err != nil {
		return nil, err
	}
	return addons, nil
}

func (i interactor) GetAddonById(service *domain.Service, id uuid.UUID) (*domain.Addon, *e.Error) {
	addon, err := i.addonRepository.GetAddonById(service, id)
	if err != nil {
		return nil, err.Append(domain.ErrAddonNotFound)
	}
	return addon, nil
}

func (i interactor) createPostgresAddon(service *domain.Service, r *request.AddonCreationRequest) (*domain.Addon, *e.Error) {
	environment, err := i.environmentRepository.GetEnvironments(service.ProjectID)

	if err != nil {
		return nil, err.Append(domain.ErrUnableToCreateAddon)
	}

	addonName := fmt.Sprintf("%s-postgres", service.DisplayName)
	addonSlug := slugify(addonName)
	serviceName := fmt.Sprintf("%s-service", addonSlug)
	volumeClaimName := fmt.Sprintf("%s-claim", addonSlug)
	volumeSize := clusterDomain.QuantityGigabytes
	postgresVolumeMountPath := "/var/lib/postgresql/data"

	addonToBeCreated := domain.Addon{
		DisplayName: addonName,
		Type:        domain.AddonTypePostgres,
		ServiceID:   service.ID,
	}

	addonToBeCreated.Slug = slugify(addonToBeCreated.DisplayName)
	addonToBeCreated.Initialize()

	addon, err := i.addonRepository.CreateAddon(&addonToBeCreated)

	if err != nil {
		return nil, err.Append(domain.ErrUnableToCreateAddon)
	}

	for _, env := range environment {
		secretName := fmt.Sprintf("%s-secret", addonSlug)
		volumeName := fmt.Sprintf("%d-%s-volume", env.ID.ID(), addonSlug)

		secretCreationRequest := clusterDomain.SecretCreationRequest{
			Name: secretName,
		}

		secret := clusterDomain.SecretUpdateRequest{
			Content: []clusterDomain.Environment{
				{
					Name:  "POSTGRES_DB",
					Value: addonSlug,
				},
				{
					Name:  "POSTGRES_USER",
					Value: generateRandomStrongString(32),
				},
				{
					Name:  "POSTGRES_PASSWORD",
					Value: generateRandomStrongString(32),
				},
				{
					Name:  "POSTGRES_HOST",
					Value: fmt.Sprintf("$%s_SERVICE_HOST", slugToEnvironmentVariable(serviceName)),
				},
				{
					Name:  "POSTGRES_PORT",
					Value: fmt.Sprintf("$%s_SERVICE_PORT", slugToEnvironmentVariable(serviceName)),
				},
			},
		}

		err := i.kubeClient.CreateSecret(env.Slug, secretCreationRequest)

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.UpdateSecret(env.Slug, secretCreationRequest.Name, secret)

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.AddSecretsToDeployment(env.Slug, service.Slug, secretCreationRequest.Name)

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.CreateDeployment(env.Slug, clusterDomain.DeploymentCreationRequest{
			DeploymentName: addon.Slug,
			ImageName:      "postgres:latest",
			Ports: []clusterDomain.Port{
				{
					ServicePort: 5432,
				},
			},
			Environment: nil,
			ConfigMaps:  nil,
			Secrets:     []string{secretName},
			Replicas:    1,
		})

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.CreateService(env.Slug, clusterDomain.Service{
			Name:           serviceName,
			DeploymentName: addon.Slug,
			PortName:       "pgsql",
			Protocol:       "TCP",
			Port:           5432,
			TargetPort:     5432,
		})

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.CreatePersistentVolume(env.Slug, clusterDomain.PersistentVolume{
			Name:        volumeName,
			StorageSize: volumeSize,
			MountPath:   postgresVolumeMountPath,
		})

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.CreatePersistentVolumeClaim(env.Slug, clusterDomain.PersistentVolumeClaim{
			Name:        volumeClaimName,
			StorageSize: volumeSize,
		})

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}

		err = i.kubeClient.AddVolumeToDeployment(env.Slug, addon.Slug, clusterDomain.VolumeDeployment{
			Name:           volumeName,
			ClaimName:      volumeClaimName,
			DeploymentName: addon.Slug,
			MountPath:      postgresVolumeMountPath,
			StorageSize:    volumeSize,
			SubPath:        "postgres",
		})

		if err != nil {
			return nil, err.Append(domain.ErrUnableToCreateAddon)
		}
	}
	return addon, nil
}

func (i interactor) CreateAddon(service *domain.Service, r *request.AddonCreationRequest) (*domain.Addon, *e.Error) {
	foundAddons, err := i.addonRepository.GetAddons(service)

	if err != nil {
		return nil, err.Append(domain.ErrUnableToCreateAddon)
	}

	if len(foundAddons) > 0 && r.Type == domain.AddonTypePostgres {
		return nil, e.Wrap(domain.ErrAServiceCanOnlyHaveOnePostgresAddon)
	}

	if r.Type == domain.AddonTypePostgres {
		return i.createPostgresAddon(service, r)
	}

	return nil, err.Append(domain.ErrUnknownAddonType)
}
