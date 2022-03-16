package ovh

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/ovhDomain"
	"github.com/ovh/go-ovh/ovh"
)

type OVHClient struct {
	Config config.OVHConfig
	client *ovh.Client
}

func (c *OVHClient) GetDomains() ([]ovhDomain.DomainName, *e.Error) {
	var domains []ovhDomain.DomainName

	err := c.client.Get("/domain", &domains)

	if err != nil {
		return nil, e.Wrap(err)
	}

	return domains, nil
}

func NewOVHClient(config config.OVHConfig) *OVHClient {
	client, _ := ovh.NewClient(
		config.Endpoint,
		config.ApplicationKey,
		config.ApplicationSecret,
		config.ConsumerKey,
	)

	return &OVHClient{
		Config: config,
		client: client,
	}
}
