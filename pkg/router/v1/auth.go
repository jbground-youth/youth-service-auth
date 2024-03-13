package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/youth-service/auth/pkg/service"
	"net/http"
)

func Login(c *gin.Context) {

	u, err := service.NewUser(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "uid is not valid : "+err.Error())
		return
	}

	if !u.CheckPassword() {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	p := service.NewPassport(u)

	err = p.GenerateToken()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = p.AddPassport()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	tokens := map[string]string{
		"access_token":  p.AccessToken,
		"refresh_token": p.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func Reissue(c *gin.Context) {

	//token data 추출
	param, err := service.MakeParamByBody(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//refresh token 검증
	if param["refresh_token"] == "" {
		c.JSON(http.StatusUnprocessableEntity, "cannot found refresh token")
		return
	}
	refreshToken := param["refresh_token"]
	token, err := service.ParseRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, "refresh expired")
		return
	}

	p, err := service.MakePassport(claims)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//delete refresh token & reissue access, refresh token

	//Delete the previous Refresh Token
	if err = p.DeleteRefreshToken(); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	//Create new pairs of refresh and access tokens
	if err = p.GenerateToken(); err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	//save the tokens metadata to redis
	if err = p.AddPassport(); err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	result := map[string]string{
		"access_token":  p.AccessToken,
		"refresh_token": p.RefreshToken,
	}
	c.JSON(http.StatusCreated, result)

}

// Logout /logout 로그아웃 시 토큰은 삭제 못하니 토큰 메타데이터 삭제
func Logout(c *gin.Context) {

	header, err := service.MakeParamByHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	token, err := service.ParseAccessToken(header["access_token"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	p, err := service.MakePassport(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
	}

	err = p.DeletePassportData()
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}
