package templates

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

var ContainerTemplate = domain.Template{
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
