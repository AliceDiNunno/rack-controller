package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/AliceDiNunno/rack-controller/src/security/crypto"
	"github.com/google/uuid"
)

func (i interactor) GetUserById(id uuid.UUID) (*userDomain.User, *e.Error) {
	user, err := i.userRepository.GetUserById(id)
	if err != nil {
		return nil, err.Append(clusterDomain.ErrUserNotFound)
	}
	return user, nil
}

func (i interactor) CreateUser(user *userDomain.UserCreationRequest) *e.Error {
	//TODO create initial user should call this function

	return e.Wrap(domain.ErrNotImplemented)
}

func (i interactor) CreateInitialUser(initialUser *config.InitialUserConfig) *e.Error {
	//TODO: check if there are no admins instead of just checking if there are no users
	if !i.userRepository.IsEmpty() {
		return e.Wrap(clusterDomain.ErrCannotCreateInitialUserIfUserTableNotEmpty)
	}

	hash, stderr := crypto.HashAndSalt(initialUser.Password)

	if stderr != nil {
		return e.Wrap(stderr)
	}

	userToCreate := &userDomain.User{
		Mail:     initialUser.Mail,
		Password: hash,
	}

	userToCreate.Initialize()

	err := i.userRepository.CreateUser(userToCreate)

	if err != nil {
		return err
	}

	hashedToken, stderr := crypto.HashAndSalt(initialUser.AccessToken)

	if stderr != nil {
		return e.Wrap(stderr)
	}

	tokenToCreate := &userDomain.AccessToken{
		User:              userToCreate,
		Token:             hashedToken,
		IsPersonnalAccess: true,
	}

	tokenToCreate.Initialize()

	err = i.userTokenRepository.CreateToken(tokenToCreate)

	if err != nil {
		return err
	}

	return nil
}
