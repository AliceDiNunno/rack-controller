package domain

import "errors"

var (
	ErrUserIsNil   = errors.New("user is nil")
	ErrIdCantBeNil = errors.New("id can't be nil")

	ErrUnableToGetProjects              = errors.New("unable to get projects")
	ErrProjectNameIsEmpty               = errors.New("project name is empty")
	ErrProjectNotFound                  = errors.New("project not found")
	ErrProjectWithSameNameAlreadyExists = errors.New("project with same name already exists")
	ErrUserNotOwnerOfProject            = errors.New("user not owner of project")
	ErrUnableToCreateProject            = errors.New("unable to create project")
)
