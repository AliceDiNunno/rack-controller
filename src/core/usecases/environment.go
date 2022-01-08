package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
)

func (i interactor) GetEnvironments(project *domain.Project) ([]domain.Environment, *e.Error) {
	environments, err := i.environmentRepository.GetEnvironments(project.ID)

	if err != nil {
		return nil, err
	}

	return environments, nil
}
