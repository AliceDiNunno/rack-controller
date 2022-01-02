package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/security/crypto"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

func (i interactor) CreateAuthToken(request domain.AccessTokenRequest) (string, *e.Error) {
	user, err := i.userRepo.FindByMail(request.Mail)
	if err != nil {
		return "", e.Wrap(domain.ErrUserNotFound)
	}

	samePassword, stderr := crypto.ComparePasswords(user.Password, request.Password)

	if stderr != nil || samePassword == false {
		return "", e.Wrap(stderr).Append(domain.ErrUserNotFound)
	}

	token := &domain.AccessToken{
		User:              user,
		IsPersonnalAccess: false,
	}

	token.Initialize()

	err = i.userTokenRepo.CreateToken(token)

	if err != nil {
		return "", err.Append(domain.ErrUserTokenCreation)
	}

	return token.Token, nil
}

func (i interactor) CreateJwtToken(request domain.JwtTokenRequest) (string, *e.Error) {
	accessToken, err := i.userTokenRepo.FindByToken(request.UserAccessToken)

	if err != nil {
		return "", e.Wrap(domain.ErrUserTokenNotFound)
	}

	if !accessToken.IsPersonnalAccess && len(accessToken.JwtGenerated) > 0 {
		return "", e.Wrap(domain.ErrJwtTokenAlreadyConsumed)
	}

	valid := accessToken.Valid()
	if !valid {
		return "", e.Wrap(domain.ErrUserTokenIsNotValid)
	}

	hostname, stderr := os.Hostname()
	if stderr != nil {
		hostname = "localhost"
	}

	var expirationDate time.Time
	if accessToken.IsPersonnalAccess {
		expirationDate = time.Now().Add(time.Hour * 24 * 7 * 30)
	} else {
		expirationDate = time.Now().Add(time.Hour * 24)
	}

	issuingDate := time.Now()
	payload := domain.JwtTokenPayload{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: expirationDate.Unix(),
			Id:        "",
			IssuedAt:  issuingDate.Unix(),
			//TODO: get a better issuer value
			Issuer:    hostname,
			NotBefore: 0,
			Subject:   accessToken.User.Mail,
		},

		UserID: accessToken.User.ID.String(),
	}

	payload.Initialize()

	mySigningKey := []byte("AllYourBase")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, stderr := token.SignedString(mySigningKey)

	if stderr != nil {
		return "", e.Wrap(stderr).Append(domain.ErrUserTokenCreation)
	}

	tkn := strings.Split(signedToken, ".")

	signature := &domain.JwtSignature{
		IssuedAt:  issuingDate,
		Token:     accessToken,
		Signature: tkn[2],
	}
	signature.Initialize()

	err = i.jwtSignature.SaveSignature(signature)

	if err != nil {
		return "", err.Append(domain.ErrUserTokenCreation)
	}

	return signedToken, nil
}

func (i interactor) CheckJwtToken(tokenStr string) (*domain.JwtTokenPayload, *e.Error) {
	token, stderr := jwt.ParseWithClaims(tokenStr, &domain.JwtTokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if stderr != nil {
		return nil, e.Wrap(domain.ErrJwtTokenCanNotBeParsed)
	}

	if !i.jwtSignature.CheckIfSignatureExists(token.Signature) {
		return nil, e.Wrap(domain.ErrJwtTokenNotTrusted)
	}

	claims := token.Claims.(*domain.JwtTokenPayload)

	if claims == nil {
		return nil, e.Wrap(domain.ErrJwtTokenClaimsInvalid)
	}

	stderr = claims.Valid()
	if stderr != nil {
		return nil, e.Wrap(stderr).Append(domain.ErrJwtTokenInvalid)
	}

	return claims, nil
}
