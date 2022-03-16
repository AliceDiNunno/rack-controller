package ip

import (
	"encoding/json"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"net/http"
)

type IpCollector struct {
}

func (i IpCollector) GetIP(ip string) (*domain.IpInformation, *e.Error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, e.Wrap(err)
	}
	defer resp.Body.Close()

	var ipInfo domain.IpInformation
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return nil, e.Wrap(err)
	}
	return &ipInfo, nil
}

func NewIPCollector() *IpCollector {
	return &IpCollector{}
}
