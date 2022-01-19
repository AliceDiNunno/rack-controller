package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases/kubernetes"
	"github.com/google/uuid"
)

type Logger interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type UserRepository interface {
	IsEmpty() bool
	CreateUser(user *userDomain.User) *e.Error
	GetUserByMail(mail string) (*userDomain.User, *e.Error)
	GetUserById(id uuid.UUID) (*userDomain.User, *e.Error)
}

type UserTokenRepository interface {
	CreateToken(token *userDomain.AccessToken) *e.Error
	FindByToken(token string) (*userDomain.AccessToken, *e.Error)
}

type JwtSignatureRepository interface {
	SaveSignature(signature *userDomain.JwtSignature) *e.Error
	CheckIfSignatureExists(signature string) bool
}

type ProjectRepository interface {
	GetProjectsByUserId(userId uuid.UUID) ([]domain.Project, *e.Error)
	GetProjectByName(name string) (*domain.Project, *e.Error)
	GetProjectByID(ID uuid.UUID) (*domain.Project, *e.Error)
	GetProjectBySlug(slug string) (*domain.Project, *e.Error)
	CreateProject(project domain.Project) *e.Error
}

type EnvironmentRepository interface {
	GetEnvironments(projectID uuid.UUID) ([]domain.Environment, *e.Error)
	GetEnvironmentByName(projectID uuid.UUID, name string) (*domain.Environment, *e.Error)
	GetEnvironmentByID(projectID uuid.UUID, ID uuid.UUID) (*domain.Environment, *e.Error)
	CreateEnvironment(environment *domain.Environment) *e.Error
}

type ServiceRepository interface {
	GetServices(ID uuid.UUID) ([]domain.Service, *e.Error)
	GetServiceByName(ID uuid.UUID, name string) (*domain.Service, *e.Error)
	CreateService(s *domain.Service) *e.Error
	UpdateService(s *domain.Service) *e.Error
}

type ConfigRepository interface {
	GetConfigByObjectID(ID uuid.UUID) ([]clusterDomain.Environment, *e.Error)
	SetConfig(ID uuid.UUID, config []clusterDomain.Environment) *e.Error
}

type EventDispatcher interface {
	Dispatch(event string, payload interface{})
	RegisterForEvent(event string, callback func(interface{}))
}

type interactor struct {
	userRepository         UserRepository
	userTokenRepository    UserTokenRepository
	jwtSignatureRepository JwtSignatureRepository
	projectRepository      ProjectRepository
	environmentRepository  EnvironmentRepository
	serviceRepository      ServiceRepository
	configRepository       ConfigRepository
	dispatcher             EventDispatcher
	kubeClient             kubernetes.Kubernetes
}

func (i interactor) GetProjectConfig(project *domain.Project) ([]clusterDomain.Environment, *e.Error) {
	if project == nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	config, err := i.configRepository.GetConfigByObjectID(project.ID)

	if err != nil {
		return nil, err.Append(domain.UnableToGetConfig)
	}

	return config, nil
}

func (i interactor) UpdateProjectConfig(project *domain.Project, envVariables []clusterDomain.Environment) *e.Error {
	//TODO implement me
	panic("implement me")
}

func (i interactor) GetEnvironmentConfig(env *domain.Environment) ([]clusterDomain.Environment, *e.Error) {
	if env == nil {
		return nil, e.Wrap(domain.ErrEnvironmentNotFound)
	}

	config, err := i.configRepository.GetConfigByObjectID(env.ID)

	if err != nil {
		return nil, err.Append(domain.UnableToGetConfig)
	}

	return config, nil
}

func (i interactor) UpdateEnvironmentConfig(env *domain.Environment, envVariables []clusterDomain.Environment) *e.Error {
	
}

func (i interactor) GetServiceConfig(service *domain.Service) ([]clusterDomain.Environment, *e.Error) {
	//TODO implement me
	panic("implement me")
}

func (i interactor) UpdateServiceConfig(service *domain.Service, envVariables []clusterDomain.Environment) *e.Error {
	//TODO implement me
	panic("implement me")
}

func NewInteractor(u UserRepository, ut UserTokenRepository, js JwtSignatureRepository,
	repo ProjectRepository, env EnvironmentRepository, s ServiceRepository, c ConfigRepository,
	kube kubernetes.Kubernetes, ed EventDispatcher) interactor {
	return interactor{
		userRepository:         u,
		userTokenRepository:    ut,
		jwtSignatureRepository: js,
		projectRepository:      repo,
		environmentRepository:  env,
		serviceRepository:      s,
		configRepository:       c,
		dispatcher:             ed,
		kubeClient:             kube,
	}
}
