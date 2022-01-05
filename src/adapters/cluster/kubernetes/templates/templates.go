package templates

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

var containerTemplate = domain.Template{
	Id:          0,
	Name:        "Container",
	Description: "Deploys a new container",
	Requirements: domain.TemplateRequirements{
		ImageName:   true,
		Environment: true,
		Ports:       true,
	},
	Exec: ExecBasicContainerTemplate,
}

func AvailableTemplateCategories() []domain.TemplateCategory {
	return []domain.TemplateCategory{
		{
			Name: "Apps",
			Templates: []domain.Template{
				containerTemplate,
			},
		},
	}
}

func AvailableTemplates() []domain.Template {
	return []domain.Template{
		containerTemplate,
	}
}
