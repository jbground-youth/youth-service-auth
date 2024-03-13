package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/youth-service/auth/pkg/msg"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}

		//bearToken := c.Request.Header.Get("Authorization")

		code := msg.SUCCESS

		if code != msg.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": msg.GetMsg(code),
				"data":    data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
