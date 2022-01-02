package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"gorm.io/gorm"
)
import "github.com/google/uuid"

type userRepo struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Mail     string
	Password string
	Tokens   []*UserToken
}

func (u userRepo) IsEmpty() bool {
	var count int64

	u.db.Model(&User{}).Count(&count)

	return count == 0
}

func (u userRepo) CreateUser(user *domain.User) *e.Error {
	userToCreate := userFromDomain(user)

	result := u.db.Create(userToCreate)

	if result.Error != nil {
		return e.Wrap(result.Error)
	}

	return nil
}

func (u userRepo) FindByMail(mail string) (*domain.User, *e.Error) {
	var user User

	result := u.db.Where("mail = ?", mail).First(&user)

	if result.Error != nil {
		return nil, e.Wrap(result.Error)
	}

	return userToDomain(&user), nil
}

func userToDomain(user *User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		ID:       user.ID,
		Mail:     user.Mail,
		Password: user.Password,
	}
}

func userFromDomain(user *domain.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:       user.ID,
		Mail:     user.Mail,
		Password: user.Password,
	}
}

func NewUserRepo(db *gorm.DB) userRepo {
	return userRepo{
		db: db,
	}
}
