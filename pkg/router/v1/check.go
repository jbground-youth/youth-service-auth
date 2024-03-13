package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/youth-service/auth/pkg/service"
	"net/http"
)

func Check(c *gin.Context) {

	body, err := service.MakeParamByBody(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	param, err := service.MakeParamByHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	token, err := service.ParseAccessToken(param["access_token"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	err = service.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	p, err := service.MakePassportByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	//err := p.ParseRedis()
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	result := map[string]string{
		"userId": p.User.UserId,
		"title":  body["title"],
	}
	c.JSON(http.StatusCreated, result)
}
