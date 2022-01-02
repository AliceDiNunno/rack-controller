package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userTokenRepo struct {
	db *gorm.DB
}

type UserToken struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;primary_key"`
	Token             string
	UserID            uuid.UUID
	User              *User
	IsPersonnalAccess bool
	JwtSignature      []*JwtSignature
}

func (u userTokenRepo) CreateToken(token *domain.AccessToken) *e.Error {
	userTokenToCreate := userTokenFromDomain(token)

	result := u.db.Create(userTokenToCreate)

	if result.Error != nil {
		return e.Wrap(result.Error)
	}

	return nil
}

func (u userTokenRepo) FindByToken(token string) (*domain.AccessToken, *e.Error) {
	userToken := &UserToken{}

	result := u.db.Joins("User").Preload("JwtSignature").Where("token = ?", token).First(userToken)

	if result.Error != nil {
		return nil, e.Wrap(result.Error)
	}

	return userTokenToDomain(userToken), nil
}

func userTokenToDomain(userToken *UserToken) *domain.AccessToken {
	if userToken == nil {
		return nil
	}

	var user *domain.User = nil

	if userToken.User != nil {
		user = userToDomain(userToken.User)
	}

	return &domain.AccessToken{
		CreatedAt:         userToken.CreatedAt,
		ID:                userToken.ID,
		Token:             userToken.Token,
		User:              user,
		IsPersonnalAccess: userToken.IsPersonnalAccess,
		JwtGenerated:      jwtSignaturesToDomain(userToken.JwtSignature),
	}
}

func userTokenFromDomain(user *domain.AccessToken) *UserToken {
	return &UserToken{
		ID:                user.ID,
		Token:             user.Token,
		UserID:            user.ID,
		User:              userFromDomain(user.User),
		IsPersonnalAccess: user.IsPersonnalAccess,
	}
}

func NewUserTokenRepo(db *gorm.DB) userTokenRepo {
	return userTokenRepo{
		db: db,
	}
}
