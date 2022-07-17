package clusterDomain

import "errors"

//TODO move to userDomain
var (
	ErrCannotCreateInitialUserIfUserTableNotEmpty = errors.New("cannot create initial user if user table is not empty")

	ErrUserNotFound        = errors.New("user not found")
	ErrUserTokenCreation   = errors.New("error creating user token")
	ErrUserTokenIsNotValid = errors.New("user token is not valid")
	ErrUserTokenNotFound   = errors.New("user token not found")

	ErrJwtTokenExpired         = errors.New("token expired")
	ErrJwtTokenAlreadyConsumed = errors.New("token already consumed")
	ErrJwtTokenInvalid         = errors.New("jwt token invalid")
	ErrJwtTokenCanNotBeParsed  = errors.New("jwt token can not be parsed")
	ErrJwtTokenClaimsInvalid   = errors.New("jwt token claims invalid")
	ErrJwtTokenNotTrusted      = errors.New("jwt token not trusted")

	ErrServiceNotFound = errors.New("service not found")

	ErrNodeNotFound = errors.New("node not found")
)
