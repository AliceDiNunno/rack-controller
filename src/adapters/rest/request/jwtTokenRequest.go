package request

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
)

type JwtTokenRequest struct {
	UserAccessToken string `binding:"required"`
}

func (r *JwtTokenRequest) ToDomain() userDomain.JwtTokenRequest {
	return userDomain.JwtTokenRequest{
		UserAccessToken: r.UserAccessToken,
	}
}
