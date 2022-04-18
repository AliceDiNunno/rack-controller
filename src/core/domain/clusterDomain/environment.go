package clusterDomain

import "fmt"

type Environment struct {
	Name  string
	Value string
}

func EnvironmentListFromMap(environmentMap map[string]interface{}) []Environment {
	var environmentList []Environment
	for key, value := range environmentMap {
		environmentList = append(environmentList, Environment{Name: key, Value: fmt.Sprintf("%v", value)})
	}
	return environmentList
}
