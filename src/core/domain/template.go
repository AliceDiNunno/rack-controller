package domain

type TemplateRequirements struct {
	ImageName   bool
	Environment bool
	Ports       bool
}

type TemplateFunction func(namespace string, request DeploymentCreationRequest) interface{}

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
