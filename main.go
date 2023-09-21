package main

import (
	"better-admin-backend-service/app"
	"better-admin-backend-service/app/db"
	"better-admin-backend-service/config"
	"better-admin-backend-service/http/rest"
	"context"
	filename "github.com/keepeye/logrus-filename"
	"github.com/open-policy-agent/opa/rego"
	log "github.com/sirupsen/logrus"
)

func main() {
	setUpLogFormatter()

	err := config.InitConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	regoQuery, err := rego.New(rego.Query("data.rest.allowed"), rego.Load([]string{
		"authorization/rest/policy.rego", "authorization/rest/data.json",
	}, nil)).PrepareForEval(context.TODO())

	if err != nil {
		log.Fatal("rego error", err)
	}

	log.Fatal(app.NewApp(rest.Router{}, db.ProductionDbConnector{}, &regoQuery).Run())
}

func setUpLogFormatter() {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
	log.SetFormatter(&log.JSONFormatter{DisableHTMLEscape: true})
}
