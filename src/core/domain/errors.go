package domain

import "errors"

var (
	ErrEnvironmentNameContainsInvalidCharacters = errors.New("environment name contains invalid characters")
	ErrEnvironmentNameCannotBeEmpty             = errors.New("environment name cannot be empty")
	ErrCPULimitationCannotBeNull                = errors.New("CPU Limitation cannot be null")
	ErrMemoryLimitationCannotBeNull             = errors.New("memory Limitation cannot be null")
	ErrReplicasCannotBeNull                     = errors.New("replicas cannot be null")
	ErrPortNameCannotBeEmpty                    = errors.New("port name cannot be empty")
	ErrPortValueIsInvalid                       = errors.New("port value is invalid")
	ErrDeploymentNameIsInvalid                  = errors.New("deployment name is invalid")
)
