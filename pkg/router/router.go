package router

import (
	"github.com/gin-gonic/gin"
	"github.com/youth-service/auth/pkg/middleware"
	"github.com/youth-service/auth/pkg/router/v1"
)

func InitializeRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.JWT())

	r.POST("/auth/login", v1.Login)
	r.POST("/auth/check", v1.Check)
	r.POST("/auth/reissue", v1.Reissue)
	r.POST("/auth/logout", v1.Logout)

	return r
}
