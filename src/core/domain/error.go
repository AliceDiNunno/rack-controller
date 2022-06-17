package domain

import "errors"

var (
	ErrUserIsNil      = errors.New("user is nil")
	ErrIdCantBeNil    = errors.New("id can't be nil")
	ErrNameTooLong    = errors.New("name is too long (max 32)")
	ErrInvalidRequest = errors.New("invalid request")

	ErrUnableToGetProjects              = errors.New("unable to get projects")
	ErrProjectNameIsEmpty               = errors.New("project name is empty")
	ErrProjectNotFound                  = errors.New("project not found")
	ErrProjectWithSameNameAlreadyExists = errors.New("project with same name already exists")
	ErrUserNotOwnerOfProject            = errors.New("user not owner of project")
	ErrUnableToCreateProject            = errors.New("unable to create project")
	ErrUnableToDeleteProject            = errors.New("unable to delete project")

	ErrEnvironmentNameIsEmpty               = errors.New("environment name is empty")
	ErrEnvironmentAlreadyExistsWithThisName = errors.New("an environment already exists with this name")
	ErrEnvironmentNotFound                  = errors.New("environment not found")
	ErrUnableToGetEnvironments              = errors.New("unable to get environments")
	ErrUnableToDeleteEnvironment            = errors.New("unable to delete environment")

	ErrUnableToGetService     = errors.New("unable to get service")
	ErrServiceNameIsEmpty     = errors.New("service name is empty")
	ErrServiceNotFound        = errors.New("service not found")
	ErrUnableToRestartService = errors.New("unable to restart service")
	ErrUnableToDeleteService  = errors.New("unable to delete service")

	ErrAddonNotFound                       = errors.New("addon not found")
	ErrUnableToCreateAddon                 = errors.New("unable to create addon")
	ErrAServiceCanOnlyHaveOnePostgresAddon = errors.New("a service can only have one postgres addon")
	ErrUnknownAddonType                    = errors.New("unknown addon type")

	UnableToGetConfig    = errors.New("unable to get config")
	UnableToUpdateConfig = errors.New("unable to update config")

	ErrInstanceNotFound = errors.New("instance not found")

	ErrNotImplemented = errors.New("not implemented")
)
