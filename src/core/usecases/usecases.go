package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/google/uuid"
)

type Usecases interface {
	//Authentication
	CreateAuthToken(request userDomain.AccessTokenRequest) (string, *e.Error)
	CreateJwtToken(request userDomain.JwtTokenRequest) (string, *e.Error)
	CheckJwtToken(token string) (*userDomain.JwtTokenPayload, *e.Error)

	CreateInitialUser(user *config.InitialUserConfig) *e.Error
	CreateUser(user *userDomain.UserCreationRequest) *e.Error
	GetUserById(id uuid.UUID) (*userDomain.User, *e.Error)

	CreateProject(user *userDomain.User, project request.CreateProjectRequest) (*domain.Project, *e.Error)
	GetUserProjects(user *userDomain.User) ([]domain.Project, *e.Error)
	GetProjectByID(user *userDomain.User, id uuid.UUID) (*domain.Project, *e.Error)

	GetEnvironments(project *domain.Project) ([]domain.Environment, *e.Error)
}
