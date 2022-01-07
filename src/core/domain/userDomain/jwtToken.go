package userDomain

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type JwtTokenRequest struct {
	UserAccessToken string `json:"token"`
}

type JwtTokenPayload struct {
	jwt.StandardClaims

	UserID uuid.UUID `json:"userId"`
}

func (j *JwtTokenPayload) Initialize() {
	j.StandardClaims.Id = uuid.New().String()
}

func (j JwtTokenPayload) Valid() error {
	expiration := time.Unix(j.ExpiresAt, 0)
	if expiration.Before(time.Now()) {
		return clusterDomain.ErrJwtTokenExpired
	}

	return nil
}
