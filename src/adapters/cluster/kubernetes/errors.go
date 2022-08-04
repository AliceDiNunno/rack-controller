package kubernetes

import "errors"

//Todo: less generic error messages
var (
	ErrDeploymentNotFound = errors.New("deployment not found")
	ErrNodeNotFound       = errors.New("node not found")
	ErrPodNotFound        = errors.New("pod not found")
	ErrTemplateNotFound   = errors.New("template not found")

	ErrConfigMapNotFound       = errors.New("config map not found")
	ErrConfigMapCreationFailed = errors.New("configmap creation failed")
	ErrConfigMapDeletionFailed = errors.New("configmap deletion failed")
	ErrConfigMapUpdateFailed   = errors.New("configmap update failed")

	ErrSecretNotFound       = errors.New("secret not found")
	ErrSecretCreationFailed = errors.New("secret creation failed")
	ErrSecretDeletionFailed = errors.New("secret deletion failed")
	ErrSecretUpdateFailed   = errors.New("secret update failed")

	ErrUnableToCreateRessource = errors.New("unable to create resource")
	ErrUnableToGetRessource    = errors.New("unable to get ressource")
	ErrUnableToDeleteRessource = errors.New("unable to delete ressource")
	ErrUnableToDeployApp       = errors.New("unable to deploy app")
	ErrUnableToUpdateApp       = errors.New("unable to update app")

	ErrPersistentVolumeClaimNotFound = errors.New("persistent volume claim not found")
	ErrPersistentVolumeNotFound      = errors.New("persistent volume not found")

	ErrUnableToCreateService = errors.New("unable to create service")
)
