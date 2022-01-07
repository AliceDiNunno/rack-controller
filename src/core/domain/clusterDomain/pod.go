package clusterDomain

import "time"

type Pod struct {
	Id           string
	Name         string
	CreationDate time.Time
	Status       string
	Image        PodImage
	ReadyProbe   *PodProbe
	NodeName     string
	QoS          string
	Condition    PodCondition
	RestartCount int
	InternalIp   string
}

type PodImage struct {
	Name  string
	Image string
}

type PodProbe struct {
	Path   string
	Scheme string
	Port   string
}

type PodCondition struct {
	Initialized     bool
	Ready           bool
	ContainersReady bool
	PodScheduled    bool
}
