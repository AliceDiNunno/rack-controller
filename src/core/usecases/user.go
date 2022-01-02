package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/security/crypto"
)

func (i interactor) CreateUser(user *domain.UserCreationRequest) *e.Error {
	panic("implement me")
}

func (i interactor) CreateInitialUser(initialUser *config.InitialUserConfig) *e.Error {
	//TODO: check if there are no admins instead of just checking if there are no users
	if !i.userRepo.IsEmpty() {
		return e.Wrap(domain.ErrCannotCreateInitialUserIfUserTableNotEmpty)
	}

	hash, stderr := crypto.HashAndSalt(initialUser.Password)

	if stderr != nil {
		return e.Wrap(stderr)
	}

	userToCreate := &domain.User{
		Mail:     initialUser.Mail,
		Password: hash,
	}

	userToCreate.Initialize()

	err := i.userRepo.CreateUser(userToCreate)

	if err != nil {
		return err
	}

	hashedToken, stderr := crypto.HashAndSalt(initialUser.AccessToken)

	if stderr != nil {
		return e.Wrap(stderr)
	}

	tokenToCreate := &domain.AccessToken{
		User:              userToCreate,
		Token:             hashedToken,
		IsPersonnalAccess: true,
	}

	tokenToCreate.Initialize()

	err = i.userTokenRepo.CreateToken(tokenToCreate)

	if err != nil {
		return err
	}

	return nil
}
