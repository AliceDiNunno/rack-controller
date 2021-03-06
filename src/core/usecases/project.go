package usecases

import (
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/google/uuid"
)

func (i interactor) GetUserProjects(user *userDomain.User) ([]domain.Project, *e.Error) {
	if user == nil {
		return nil, e.Wrap(domain.ErrUserIsNil)
	}

	projects, err := i.projectRepository.GetProjectsByUserId(user.ID)
	if err != nil {
		return nil, e.Wrap(domain.ErrUnableToGetProjects)
	}

	return projects, nil
}

func (i interactor) GetProjectByID(user *userDomain.User, id uuid.UUID) (*domain.Project, *e.Error) {
	if id == uuid.Nil {
		return nil, e.Wrap(domain.ErrIdCantBeNil)
	}

	project, err := i.projectRepository.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	if project.UserID != user.ID {
		return nil, e.Wrap(domain.ErrUserNotOwnerOfProject)
	}

	return project, nil
}

func (i interactor) generateKubernetesCompatibleSlug(ressource string, name string) string {
	ressource = slugify(ressource)
	name = slugify(name)

	for len(ressource)+1+len(name) > 63 {
		ressource = ressource[:len(ressource)-1]
	}

	return fmt.Sprintf("%s-%s", ressource, name)
}

func (i interactor) CreateProject(user *userDomain.User, projectCreationRequest request.CreateProjectRequest) (*domain.Project, *e.Error) {
	if user == nil {
		return nil, e.Wrap(domain.ErrUserIsNil)
	}

	if projectCreationRequest.Name == "" {
		return nil, e.Wrap(domain.ErrProjectNameIsEmpty)
	}

	if len(projectCreationRequest.Name) > 32 {
		return nil, e.Wrap(domain.ErrNameTooLong)
	}

	_, err := i.projectRepository.GetProjectByName(projectCreationRequest.Name)

	if err == nil {
		return nil, e.Wrap(domain.ErrProjectWithSameNameAlreadyExists)
	}

	project := domain.Project{
		DisplayName: projectCreationRequest.Name,
		Slug:        i.generateKubernetesCompatibleSlug("project", projectCreationRequest.Name),
		UserID:      user.ID,
	}

	project.Initialize()

	err = i.projectRepository.CreateProject(project)

	if err != nil {
		return nil, e.Wrap(domain.ErrUnableToCreateProject)
	}

	return &project, nil
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
	if project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	err := i.configRepository.SetConfig(project.ID, envVariables)

	if err != nil {
		return err.Append(domain.UnableToUpdateConfig)
	}

	return nil
}

func (i interactor) ConfigForProject(project *domain.Project) map[string]string {
	config := map[string]string{
		"HTTP_LISTEN_ADDRESS": "0.0.0.0",
		"HTTP_LISTEN_PORT":    "80",
		"HTTP_TLS_ENABLED":    "false",

		"LOGGER_PROJECT_ID": project.ID.String(),
		"LOGGER_EVENT_KEY":  project.EventKey.String(),
	}

	if project == nil {
		return config
	}

	projectConfig, err := i.configRepository.GetConfigByObjectID(project.ID)
	if err == nil {
		for _, env := range projectConfig {
			config[env.Name] = env.Value
		}
	}

	return config
}

func (i interactor) DeleteProject(project *domain.Project) *e.Error {
	if project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	environments, err := i.GetEnvironments(project)

	if err != nil {
		return err.Append(domain.ErrUnableToDeleteProject)
	}

	for _, environment := range environments {
		err = i.DeleteEnvironment(&environment)
		if err != nil {
			return err.Append(domain.ErrUnableToDeleteProject)
		}
	}

	err = i.projectRepository.DeleteProject(project)

	if err != nil {
		return err.Append(domain.ErrUnableToDeleteProject)
	}

	return nil
}
