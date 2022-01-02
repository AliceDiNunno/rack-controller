package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
)

type Usecases interface {
	//Authentication
	CreateAuthToken(request domain.AccessTokenRequest) (string, *e.Error)
	CreateJwtToken(request domain.JwtTokenRequest) (string, *e.Error)
	CheckJwtToken(token string) (*domain.JwtTokenPayload, *e.Error)

	CreateInitialUser(user *config.InitialUserConfig) *e.Error
	CreateUser(user *domain.UserCreationRequest) *e.Error
}
