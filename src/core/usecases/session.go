package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/AliceDiNunno/rack-controller/src/security/crypto"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

func (i interactor) CreateAuthToken(request userDomain.AccessTokenRequest) (string, *e.Error) {
	user, err := i.userRepository.GetUserByMail(request.Mail)
	if err != nil {
		return "", e.Wrap(clusterDomain.ErrUserNotFound)
	}

	samePassword, stderr := crypto.ComparePasswords(user.Password, request.Password)

	if stderr != nil || samePassword == false {
		return "", e.Wrap(stderr).Append(clusterDomain.ErrUserNotFound)
	}

	token := &userDomain.AccessToken{
		User:              user,
		IsPersonnalAccess: false,
	}

	token.Initialize()

	err = i.userTokenRepository.CreateToken(token)

	if err != nil {
		return "", err.Append(clusterDomain.ErrUserTokenCreation)
	}

	return token.Token, nil
}

func (i interactor) CreateJwtToken(request userDomain.JwtTokenRequest) (string, *e.Error) {
	accessToken, err := i.userTokenRepository.FindByToken(request.UserAccessToken)

	if err != nil {
		return "", e.Wrap(clusterDomain.ErrUserTokenNotFound)
	}

	if !accessToken.IsPersonnalAccess && len(accessToken.JwtGenerated) > 0 {
		return "", e.Wrap(clusterDomain.ErrJwtTokenAlreadyConsumed)
	}

	valid := accessToken.Valid()
	if !valid {
		return "", e.Wrap(clusterDomain.ErrUserTokenIsNotValid)
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
	payload := userDomain.JwtTokenPayload{
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

		UserID: accessToken.User.ID,
	}

	payload.Initialize()

	//TODO: get a better secret value
	mySigningKey := []byte("AllYourBase")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, stderr := token.SignedString(mySigningKey)

	if stderr != nil {
		return "", e.Wrap(stderr).Append(clusterDomain.ErrUserTokenCreation)
	}

	tkn := strings.Split(signedToken, ".")

	signature := &userDomain.JwtSignature{
		IssuedAt:  issuingDate,
		Token:     accessToken,
		Signature: tkn[2],
	}
	signature.Initialize()

	err = i.jwtSignatureRepository.SaveSignature(signature)

	if err != nil {
		return "", err.Append(clusterDomain.ErrUserTokenCreation)
	}

	return signedToken, nil
}

func (i interactor) CheckJwtToken(tokenStr string) (*userDomain.JwtTokenPayload, *e.Error) {
	token, stderr := jwt.ParseWithClaims(tokenStr, &userDomain.JwtTokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if stderr != nil {
		return nil, e.Wrap(clusterDomain.ErrJwtTokenCanNotBeParsed)
	}

	if !i.jwtSignatureRepository.CheckIfSignatureExists(token.Signature) {
		return nil, e.Wrap(clusterDomain.ErrJwtTokenNotTrusted)
	}

	claims := token.Claims.(*userDomain.JwtTokenPayload)

	if claims == nil {
		return nil, e.Wrap(clusterDomain.ErrJwtTokenClaimsInvalid)
	}

	stderr = claims.Valid()
	if stderr != nil {
		return nil, e.Wrap(stderr).Append(clusterDomain.ErrJwtTokenInvalid)
	}

	return claims, nil
}
