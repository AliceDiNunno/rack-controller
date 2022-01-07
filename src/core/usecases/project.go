package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
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

func (i interactor) CreateProject(user *userDomain.User, projectCreationRequest request.CreateProjectRequest) (*domain.Project, *e.Error) {
	if user == nil {
		return nil, e.Wrap(domain.ErrUserIsNil)
	}

	if projectCreationRequest.Name == "" {
		return nil, e.Wrap(domain.ErrProjectNameIsEmpty)
	}

	_, err := i.projectRepository.GetProjectByName(projectCreationRequest.Name)

	if err == nil {
		return nil, e.Wrap(domain.ErrProjectWithSameNameAlreadyExists)
	}

	var project domain.Project

	project.DisplayName = projectCreationRequest.Name
	project.UserID = user.ID

	project.Initialize()

	err = i.projectRepository.CreateProject(project)

	if err != nil {
		return nil, e.Wrap(domain.ErrUnableToCreateProject)
	}

	return &project, nil
}
