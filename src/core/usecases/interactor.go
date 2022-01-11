package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
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
	CreateProject(project domain.Project) *e.Error
}

type EnvironmentRepository interface {
	GetEnvironments(projectID uuid.UUID) ([]domain.Environment, *e.Error)
	GetEnvironmentByName(projectID uuid.UUID, name string) (*domain.Environment, *e.Error)
	CreateEnvironment(environment *domain.Environment) *e.Error
}

type ServiceRepository interface {
	GetServices(id uuid.UUID) ([]domain.Service, *e.Error)
	GetServiceByName(id uuid.UUID, name string) (*domain.Service, *e.Error)
	CreateService(s *domain.Service) *e.Error
}

type EventDispatcher interface {
	Dispatch(event string, payload interface{})
	RegisterForEvent(event string, callback func(interface{}))
}

type Kubernetes interface {
	GetNodes() ([]clusterDomain.Node, *e.Error)
	GetNode(string) (*clusterDomain.Node, *e.Error)

	ListDeployments(namespace string) ([]clusterDomain.Deployment, *e.Error)
	GetDeployment(namespace string, deploymentName string) (*clusterDomain.Deployment, *e.Error)
	GetPodsOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Pod, *e.Error)
	GetEnvironmentOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Environment, *e.Error)
	GetPortsOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Port, *e.Error)
	GetConfigMapsOfADeployment(namespace string, name string) ([]string, *e.Error)
	GetSecretsOfADeployment(namespace string, name string) ([]string, *e.Error)
	DeleteDeployment(namespace string, name string) *e.Error
	RestartDeployment(namespace string, name string) *e.Error
	CreateDeployment(namespace string, request clusterDomain.DeploymentCreationRequest) *e.Error
	ExposePorts(namespace string, name string, ports []clusterDomain.Port) *e.Error

	GetPods(namespace string) ([]clusterDomain.Pod, *e.Error)
	GetPod(namespace string, podName string) (*clusterDomain.Pod, *e.Error)
	GetPodLogs(namespace string, podName string) (string, *e.Error)
	DeletePod(namespace string, podName string) *e.Error

	GetNamespaces() ([]string, *e.Error)

	GetConfigMapList(namespace string) ([]string, *e.Error)
	GetConfigMap(namespace string, name string) (clusterDomain.ConfigMap, *e.Error)
	CreateConfigMap(namespace string, request clusterDomain.ConfigMapCreationRequest) *e.Error
	UpdateConfigMap(namespace string, name string, request clusterDomain.ConfigMapUpdateRequest) *e.Error
	DeleteConfigMap(namespace string, name string) *e.Error

	GetSecretsList(namespace string) ([]string, *e.Error)
	GetSecret(namespace string, name string) (clusterDomain.Secret, *e.Error)
	CreateSecret(namespace string, request clusterDomain.SecretCreationRequest) *e.Error
	UpdateSecret(namespace string, name string, request clusterDomain.SecretUpdateRequest) *e.Error
	DeleteSecret(namespace string, name string) *e.Error
}

type interactor struct {
	userRepository         UserRepository
	userTokenRepository    UserTokenRepository
	jwtSignatureRepository JwtSignatureRepository
	projectRepository      ProjectRepository
	environmentRepository  EnvironmentRepository
	serviceRepository      ServiceRepository
	kubernetes             Kubernetes
	dispatcher             EventDispatcher
}

func NewInteractor(u UserRepository, ut UserTokenRepository, js JwtSignatureRepository,
	repo ProjectRepository, env EnvironmentRepository, s ServiceRepository,
	k8s Kubernetes, ed EventDispatcher) interactor {
	return interactor{
		userRepository:         u,
		userTokenRepository:    ut,
		jwtSignatureRepository: js,
		projectRepository:      repo,
		environmentRepository:  env,
		serviceRepository:      s,
		kubernetes:             k8s,
		dispatcher:             ed,
	}
}
