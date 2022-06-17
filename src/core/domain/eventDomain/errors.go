package domain

import "errors"

var (
	ErrGroupNotFound        = errors.New("group not found")
	ErrUnableToDeleteObject = errors.New("unable to delete object")
	ErrUnableToFindEvents   = errors.New("unable to find events")
	ErrInvalidProjectKey    = errors.New("invalid project key")
)
