package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func BindParameterByHeader(r *http.Request) (map[string]string, error) {
	param := map[string]string{}

	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		param["access_token"] = strArr[1]
		return param, nil
	}
	return nil, errors.New("cannot found access token")
}

func BindParameterByBody(c *gin.Context) (map[string]string, error) {
	param := map[string]string{}
	if err := c.ShouldBindJSON(&param); err != nil {
		return nil, err
	}
	return param, nil
}
