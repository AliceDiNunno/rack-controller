package request

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

type JwtTokenRequest struct {
	UserAccessToken string `binding:"required"`
}

func (r *JwtTokenRequest) ToDomain() domain.JwtTokenRequest {
	return domain.JwtTokenRequest{
		UserAccessToken: r.UserAccessToken,
	}
}
