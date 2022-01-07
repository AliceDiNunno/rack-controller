package utils

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	v16 "k8s.io/client-go/applyconfigurations/core/v1"
)

func EnvironmentVariableFromDomain(envVariable clusterDomain.Environment) v16.EnvVarApplyConfiguration {
	return v16.EnvVarApplyConfiguration{
		Name:      &envVariable.Name,
		Value:     &envVariable.Value,
		ValueFrom: nil,
	}
}

func EnvironmentVariablesFromDomain(envVariables []clusterDomain.Environment) []v16.EnvVarApplyConfiguration {
	var envVariablesList []v16.EnvVarApplyConfiguration

	for _, envVariable := range envVariables {
		envVariablesList = append(envVariablesList, EnvironmentVariableFromDomain(envVariable))
	}

	return envVariablesList
}
