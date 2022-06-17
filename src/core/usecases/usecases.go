package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	eventDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/ovhDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/google/uuid"
)

type Usecases interface {
	//Authentication
	CreateAuthToken(request userDomain.AccessTokenRequest) (string, *e.Error)
	CreateJwtToken(request userDomain.JwtTokenRequest) (string, *e.Error)
	CheckJwtToken(token string) (*userDomain.JwtTokenPayload, *e.Error)

	//Users
	CreateInitialUser(user *config.InitialUserConfig) *e.Error
	CreateUser(user *userDomain.UserCreationRequest) *e.Error
	GetUserById(id uuid.UUID) (*userDomain.User, *e.Error)

	//Projects
	CreateProject(user *userDomain.User, project request.CreateProjectRequest) (*domain.Project, *e.Error)
	GetUserProjects(user *userDomain.User) ([]domain.Project, *e.Error)
	GetProjectByID(user *userDomain.User, id uuid.UUID) (*domain.Project, *e.Error)
	GetProjectConfig(project *domain.Project) ([]clusterDomain.Environment, *e.Error)
	UpdateProjectConfig(project *domain.Project, envVariables []clusterDomain.Environment) *e.Error
	DeleteProject(project *domain.Project) *e.Error

	//Environments
	CreateEnvironment(project *domain.Project, r *request.EnvironmentCreationRequest) *e.Error
	GetEnvironments(project *domain.Project) ([]domain.Environment, *e.Error)
	GetEnvironmentByID(project *domain.Project, id uuid.UUID) (*domain.Environment, *e.Error)
	GetEnvironmentConfig(env *domain.Environment) ([]clusterDomain.Environment, *e.Error)
	UpdateEnvironmentConfig(env *domain.Environment, envVariables []clusterDomain.Environment) *e.Error
	DeleteEnvironment(env *domain.Environment) *e.Error

	//Services
	CreateService(project *domain.Project, r *request.ServiceCreationRequest) *e.Error
	GetServices(project *domain.Project) ([]domain.Service, *e.Error)
	GetServiceOfEnvironment(service *domain.Service, environment *domain.Environment) (*domain.ServiceDetail, *e.Error)
	GetServiceById(project *domain.Project, id uuid.UUID) (*domain.Service, *e.Error)
	GetServiceConfig(service *domain.Service) ([]clusterDomain.Environment, *e.Error)
	UpdateServiceConfig(service *domain.Service, envVariables []clusterDomain.Environment) *e.Error
	RestartService(service *domain.Service) *e.Error
	DeleteService(service *domain.Service) *e.Error

	//Events
	PushNewEvent(id uuid.UUID, request *request.ItemCreationRequest) *e.Error

	GetProjectsEvent(user *userDomain.User, project *domain.Project) ([]eventDomain.Event, *e.Error)
	FetchGroupingIdContent(project *domain.Project, groupingId string) (*eventDomain.Event, *e.Error)
	FetchGroupingIdOccurrences(project *domain.Project, groupingId string) ([]eventDomain.Event, *e.Error)
	FetchGroupOccurrence(project *domain.Project, groupingId string, occurrence string) (*eventDomain.Event, *e.Error)
	FetchProjectVersions(project *domain.Project) ([]eventDomain.Event, *e.Error)
	FetchProjectEnvironments(project *domain.Project) ([]eventDomain.Event, *e.Error)
	FetchProjectServers(project *domain.Project) ([]eventDomain.Event, *e.Error)

	//Instances
	GetServiceInstances(service *domain.Service, environments *domain.Environment) ([]clusterDomain.Pod, *e.Error)
	GetSpecificNodeInstances(id string) ([]clusterDomain.Pod, *e.Error)
	GetInstanceLogs(service *domain.Service, environment *domain.Environment, instance *clusterDomain.Pod) (string, *e.Error)
	GetInstanceByName(service *domain.Service, environment *domain.Environment, name string) (*clusterDomain.Pod, *e.Error)
	DeleteInstance(service *domain.Service, environment *domain.Environment, instance *clusterDomain.Pod) *e.Error

	//Addons
	GetAddons(service *domain.Service) ([]domain.Addon, *e.Error)
	GetAddonById(service *domain.Service, id uuid.UUID) (*domain.Addon, *e.Error)
	CreateAddon(service *domain.Service, r *request.AddonCreationRequest) (*domain.Addon, *e.Error)

	//Nodes
	GetNodes() ([]clusterDomain.Node, *e.Error)

	//Domain names and ingress
	GetSpecificNode(id string) (*clusterDomain.Node, *e.Error)

	GetDomainNames() ([]ovhDomain.DomainName, *e.Error)

	//IPs
	/*GetIPs(ip ...string) ([]string, *e.Error)
	GetLocalIP() (string, *e.Error)*/
}
