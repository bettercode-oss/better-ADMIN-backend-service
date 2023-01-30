package main

import (
	"better-admin-backend-service/app"
	"better-admin-backend-service/app/db"
	"better-admin-backend-service/config"
	"better-admin-backend-service/http/rest"
	filename "github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
)

func main() {
	setUpLogFormatter()

	err := config.InitConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.NewApp(rest.Router{}, db.ProductionDbConnector{}).Run())
}

func setUpLogFormatter() {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
	log.SetFormatter(&log.JSONFormatter{DisableHTMLEscape: true})
}
