package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youth-service/auth/pkg/logging"
	"github.com/youth-service/auth/pkg/router"
	"github.com/youth-service/auth/pkg/setting"
	"github.com/youth-service/auth/pkg/store"
	"github.com/youth-service/auth/pkg/store/models"
	"log"
	"net/http"
)

func Run() {

	models.Setup()
	log.Printf("[INFO] Connected mariadb database : %s", setting.DatabaseSetting.Host)

	store.RunRedis()
	log.Printf("[INFO] Connected redis database : %s", setting.RedisSetting.Host)

	logging.Setup()

	gin.SetMode(setting.ServerSetting.RunMode)

	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        router.InitializeRouter(),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[INFO] start http server listening %s", endPoint)
	log.Fatal(server.ListenAndServe())
}
