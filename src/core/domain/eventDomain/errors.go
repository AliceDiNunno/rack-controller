package domain

import "errors"

var (
	ErrFailedToGetUser                    = errors.New("failed to fetch user")
	ErrProjectNotFound                    = errors.New("project not found")
	ErrGroupNotFound                      = errors.New("group not found")
	ErrProjectAlreadyExistingWithThisName = errors.New("a project already exists with this name")
	ErrUnableToDeleteObject               = errors.New("unable to delete object")
	ErrUnknownDBError                     = errors.New("an unknown database error has happened")
)
