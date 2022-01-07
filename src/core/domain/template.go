package domain

import "github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"

type TemplateRequirements struct {
	ImageName   bool
	Environment bool
	Ports       bool
}

type TemplateFunction func(namespace string, request clusterDomain.DeploymentCreationRequest) interface{}

type Template struct {
	Id           int
	Name         string
	Description  string
	Requirements TemplateRequirements

	Exec TemplateFunction `json:"-"`
}

type TemplateCategory struct {
	Name      string
	Templates []Template
}
