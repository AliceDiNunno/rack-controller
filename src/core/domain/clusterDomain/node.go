package clusterDomain

import (
	"time"
)

type Node struct {
	Id                string
	Name              string
	CreationDate      time.Time
	Ip                string
	Hardware          NodeHardware
	AvailableHardware NodeHardware
	OperatingSystem   NodeOperatingSystem
	Condition         NodeCondition
}

type NodeTaint struct {
	Key    string
	Effect string
	Since  time.Time
}

type NodeCondition struct {
	NetworkUnavailable bool
	DiskPressure       bool
	PidPressure        bool
	Ready              bool

	Taints []NodeTaint
}

type NodeOperatingSystem struct {
	OSType         string
	OSName         string
	OSArchitecture string
	KernelVersion  string
}

type NodeHardware struct {
	Cores   int64
	Storage int64
	Memory  int64
}
