package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	logDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/ovhDomain"
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
	GetProjectByIDAndKey(id uuid.UUID, key uuid.UUID) (*domain.Project, *e.Error)
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
	CreateOrUpdateService(s *domain.Service) *e.Error
	UpdateService(s *domain.Service) *e.Error
	GetServiceById(projectID uuid.UUID, ID uuid.UUID) (*domain.Service, *e.Error)
	DeleteService(service *domain.Service) *e.Error
}

type ConfigRepository interface {
	GetConfigByObjectID(ID uuid.UUID) ([]clusterDomain.Environment, *e.Error)
	SetConfig(ID uuid.UUID, config []clusterDomain.Environment) *e.Error
}

type EventDispatcher interface {
	Dispatch(event string, payload interface{})
	RegisterForEvent(event string, callback func(interface{}))
}

type EventRepository interface {
	AddLog(log *logDomain.LogEntry) *e.Error

	ProjectVersions(project *domain.Project) ([]logDomain.LogEntry, *e.Error)
	ProjectEnvironments(project *domain.Project) ([]logDomain.LogEntry, *e.Error)
	ProjectServers(project *domain.Project) ([]logDomain.LogEntry, *e.Error)
	ProjectGroupingIds(project *domain.Project) ([]logDomain.LogEntry, *e.Error)
	IsGroupExist(project *domain.Project, groupingId string) bool

	FindLastEntryForGroup(project *domain.Project, groupingId string) (*logDomain.LogEntry, *e.Error)
	FindGroupOccurrences(project *domain.Project, groupingId string) ([]logDomain.LogEntry, *e.Error)
	FindGroupOccurrence(project *domain.Project, groupingId string, occurenceId string) (*logDomain.LogEntry, *e.Error)
}

type IpInformationCollector interface {
	GetIP(ip string) (*domain.IpInformation, *e.Error)
}

type OvhClient interface {
	GetDomains() ([]ovhDomain.DomainName, *e.Error)
}

type interactor struct {
	userRepository         UserRepository
	userTokenRepository    UserTokenRepository
	jwtSignatureRepository JwtSignatureRepository
	projectRepository      ProjectRepository
	environmentRepository  EnvironmentRepository
	serviceRepository      ServiceRepository
	configRepository       ConfigRepository
	logCollection          EventRepository
	dispatcher             EventDispatcher
	kubeClient             kubernetes.Kubernetes
	ipCollector            IpInformationCollector
	ovhClient              OvhClient
}

func NewInteractor(u UserRepository, ut UserTokenRepository, js JwtSignatureRepository,
	repo ProjectRepository, env EnvironmentRepository, s ServiceRepository, c ConfigRepository,
	lC EventRepository,
	kube kubernetes.Kubernetes, ed EventDispatcher, iic IpInformationCollector, oc OvhClient) interactor {
	return interactor{
		userRepository:         u,
		userTokenRepository:    ut,
		jwtSignatureRepository: js,
		projectRepository:      repo,
		environmentRepository:  env,
		serviceRepository:      s,
		configRepository:       c,
		logCollection:          lC,
		dispatcher:             ed,
		kubeClient:             kube,
		ipCollector:            iic,
		ovhClient:              oc,
	}
}
