package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (rH RoutesHandler) verifyAuthenticationMiddleware(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")

	if authorizationHeader == "" {
		rH.handleError(c, e.Wrap(ErrAuthorizationHeaderMissing))
		return
	}

	payload, err := rH.usecases.CheckJwtToken(authorizationHeader)

	if err != nil {
		rH.handleError(c, err.Append(ErrInvalidAuthorizationHeader))
		return
	}

	c.Set("userID", payload.UserID)
}

func (rH RoutesHandler) createAuthTokenHandler(c *gin.Context) {
	var tokenRequest request.AccessTokenRequest

	if stderr := c.ShouldBindJSON(&tokenRequest); stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	token, err := rH.usecases.CreateAuthToken(tokenRequest.ToDomain())

	if err != nil {
		rH.handleError(c, err)
		return
	}

	tokenResponse := response.AccessTokenResponse{
		AccessToken: token,
	}

	c.JSON(http.StatusCreated, success(tokenResponse))
}

func (rH RoutesHandler) createJwtTokenHandler(c *gin.Context) {
	var jwtRequest request.JwtTokenRequest

	if stderr := c.ShouldBindJSON(&jwtRequest); stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	token, err := rH.usecases.CreateJwtToken(jwtRequest.ToDomain())

	if err != nil {
		rH.handleError(c, err)
		return
	}

	jwtResponse := response.JwtTokenResponse{
		JwtToken: token,
	}

	c.JSON(http.StatusCreated, success(jwtResponse))
}

func (rH RoutesHandler) deleteJwtTokenHandler(c *gin.Context) {
	//TODO: implement
}
