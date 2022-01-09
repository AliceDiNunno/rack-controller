package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/davecgh/go-spew/spew"
)

func (i interactor) GetEnvironments(project *domain.Project) ([]domain.Environment, *e.Error) {
	environments, err := i.environmentRepository.GetEnvironments(project.ID)

	if err != nil {
		return nil, err
	}

	return environments, nil
}

func (i interactor) CreateEnvironment(project *domain.Project, r *request.EnvironmentCreationRequest) *e.Error {
	if project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	if r == nil {
		return e.Wrap(domain.ErrInvalidRequest)
	}

	if r.Name == "" {
		return e.Wrap(domain.ErrEnvironmentNameIsEmpty)
	}

	env, err := i.environmentRepository.GetEnvironmentByName(project.ID, r.Name)

	spew.Dump(env, err)

	if err == nil && env != nil {
		return e.Wrap(domain.ErrEnvironmentAlreadyExistsWithThisName)
	}

	environment := domain.Environment{
		DisplayName: r.Name,
		ProjectId:   project.ID,
	}

	environment.Initialize()

	return i.environmentRepository.CreateEnvironment(&environment)
}
