package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
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

	env, stderr := i.environmentRepository.GetEnvironmentByName(project.ID, r.Name)

	spew.Dump(env, stderr)

	if stderr == nil && env != nil {
		return e.Wrap(domain.ErrEnvironmentAlreadyExistsWithThisName)
	}

	environment := domain.Environment{
		DisplayName: r.Name,
		ProjectId:   project.ID,
		//this will appear as "project-projectname-environmentname"
		Slug: i.generateKubernetesCompatibleSlug(project.Slug, r.Name),
	}

	environment.Initialize()

	err := i.kubeClient.CreateNamespace(environment.Slug)

	if err != nil {
		return err
	}

	return i.environmentRepository.CreateEnvironment(&environment)
}

func (i interactor) GetEnvironmentByID(project *domain.Project, id uuid.UUID) (*domain.Environment, *e.Error) {
	if project == nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	if id == uuid.Nil {
		return nil, e.Wrap(domain.ErrInvalidRequest)
	}

	environment, err := i.environmentRepository.GetEnvironmentByID(project.ID, id)

	if err != nil {
		return nil, err.Append(domain.ErrEnvironmentNotFound)
	}

	return environment, nil
}

func (i interactor) EnvVariablesForEnvironment(environment *domain.Environment) map[string]string {
	if environment.DisplayName == "prod" || environment.DisplayName == "production" {
		return map[string]string{
			"ENV": "production",
		}
	} else {
		return map[string]string{
			"ENV": "development",
		}
	}
}
