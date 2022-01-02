package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type jwtSignatureRepo struct {
	db *gorm.DB
}

type JwtSignature struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Signature   string
	UserTokenID uuid.UUID
	UserToken   *UserToken
	IssuedAt    time.Time
}

func (j jwtSignatureRepo) SaveSignature(signature *domain.JwtSignature) *e.Error {
	jwtSignature := jwtSignatureFromDomain(signature)
	if err := j.db.Create(&jwtSignature).Error; err != nil {
		return e.Wrap(err)
	}
	return nil
}

func (j jwtSignatureRepo) CheckIfSignatureExists(signatureStr string) bool {
	var jwtSignature JwtSignature
	if err := j.db.Where("signature = ?", signatureStr).First(&jwtSignature).Error; err != nil {
		return false
	}
	return true
}

func jwtSignaturesToDomain(signature []*JwtSignature) []*domain.JwtSignature {
	var signaturesSlice []*domain.JwtSignature

	for _, s := range signature {
		signaturesSlice = append(signaturesSlice, jwtSignatureToDomain(s))
	}

	return signaturesSlice
}

func jwtSignatureToDomain(signature *JwtSignature) *domain.JwtSignature {
	if signature == nil {
		return nil
	}
	return &domain.JwtSignature{
		ID:        signature.ID,
		Signature: signature.Signature,
		Token:     userTokenToDomain(signature.UserToken),
		IssuedAt:  signature.IssuedAt,
	}
}

func jwtSignatureFromDomain(signature *domain.JwtSignature) *JwtSignature {
	if signature == nil {
		return nil
	}
	return &JwtSignature{
		ID:        signature.ID,
		Signature: signature.Signature,
		UserToken: userTokenFromDomain(signature.Token),
		IssuedAt:  signature.IssuedAt,
	}
}

func NewJwtSignatureRepo(db *gorm.DB) jwtSignatureRepo {
	return jwtSignatureRepo{
		db: db,
	}
}
