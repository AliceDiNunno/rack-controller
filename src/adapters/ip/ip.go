package ip

import (
	"encoding/json"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"net/http"
	"time"
)

type IpCollector struct {
	CachedLocalIP     *domain.IpInformation
	CachedLocalIPDate int64
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

func (i IpCollector) GetLocalIP() (*domain.IpInformation, *e.Error) {
	if i.CachedLocalIP != nil && i.CachedLocalIPDate > (time.Now().Unix()-60) {
		return i.CachedLocalIP, nil
	}

	ip, err := i.GetIP("")
	if err != nil {
		return nil, err
	}
	i.CachedLocalIP = ip
	i.CachedLocalIPDate = time.Now().Unix()
	return ip, nil
}

func NewIPCollector() *IpCollector {
	return &IpCollector{}
}
