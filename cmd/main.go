package main

import (
	"github.com/youth-service/auth/cmd/app"
	"github.com/youth-service/auth/pkg/setting"
)

func init() {
	setting.Setup()
}

func main() {
	app.Run()
}
