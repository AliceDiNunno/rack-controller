package request

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

type AccessTokenRequest struct {
	Mail     string `binding:"required"`
	Password string `binding:"required"`
}

func (r AccessTokenRequest) ToDomain() domain.AccessTokenRequest {
	return domain.AccessTokenRequest{
		Mail:     r.Mail,
		Password: r.Password,
	}
}
