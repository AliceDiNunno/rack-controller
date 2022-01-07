package userDomain

import (
	"github.com/google/uuid"
	"time"
)

type JwtSignature struct {
	CreatedAt time.Time
	IssuedAt  time.Time
	ID        uuid.UUID
	Token     *AccessToken
	Signature string
}

func (js *JwtSignature) Initialize() {
	js.ID = uuid.New()
}
