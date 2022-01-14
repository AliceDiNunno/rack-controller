package events

import (
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"strings"
)

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}

//TODO: registered slug database to avoid duplicates
// ProjectCreatedEvent TODO: object name should conform to dns name conventions otherwise kubernetes will not be able to create the object
func (h EventHandler) ProjectCreatedEvent(data interface{}) {
	println("project created event")

	project, ok := data.(domain.Project)

	if !ok {
		return
	}
 
	projectId := strings.Replace(project.ID.String(), "-", "", -1)
	projectName := strip(strings.ToLower(project.DisplayName))

	objectName := fmt.Sprintf("project-%s-%s", projectId, projectName)
	println(objectName)

	h.cluster.CreateNamespace(objectName)
}
