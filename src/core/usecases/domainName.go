package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/ovhDomain"
)

func (i interactor) GetDomainNames() ([]ovhDomain.DomainName, *e.Error) {
	return i.ovhClient.GetDomains()
}
