package userDomain

import "github.com/google/uuid"

type UserCreationRequest struct {
	Mail     string
	Password string
}

type User struct {
	ID       uuid.UUID
	Mail     string
	Password string `json:"-"`
}

func (u *User) Initialize() {
	u.ID = uuid.New()
}
