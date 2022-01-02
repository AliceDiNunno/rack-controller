package domain

import "time"

type Container struct {
	Name         string
	CreationDate time.Time
	Image        ContainerImage
	ReadyProbe   *ContainerProbe
}

type ContainerImage struct {
	Name  string
	Image string
}

type ContainerProbe struct {
	Path   string
	Scheme string
	Port   string
}
