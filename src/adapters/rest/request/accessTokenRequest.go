package request

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
)

type AccessTokenRequest struct {
	Mail     string `binding:"required"`
	Password string `binding:"required"`
}

func (r AccessTokenRequest) ToDomain() userDomain.AccessTokenRequest {
	return userDomain.AccessTokenRequest{
		Mail:     r.Mail,
		Password: r.Password,
	}
}
