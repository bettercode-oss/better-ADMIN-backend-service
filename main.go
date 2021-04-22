package main

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/controllers"
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/middlewares"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/keepeye/logrus-filename"
	"github.com/labstack/echo"
	echomiddleware "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	Version = "0.0.1"
	website = "https://www.bettercode.kr"
	banner  = `
  __        __   __                    ___    ___    __  ___   ____   _  __
  / /  ___  / /_ / /_ ___   ____       / _ |  / _ \  /  |/  /  /  _/  / |/ /
 / _ \/ -_)/ __// __// -_) / __/      / __ | / // / / /|_/ /  _/ /   /    / 
/_.__/\__/ \__/ \__/ \__/ /_/        /_/ |_|/____/ /_/  /_/  /___/  /_/|_/  
`
)

var (
	gormDB *gorm.DB
)

func init() {
	config.InitConfig("config/config.json")

	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
	log.SetFormatter(&log.JSONFormatter{DisableHTMLEscape: true})

	db, err := gorm.Open(sqlite.Open("account.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Database Connection Error")
	}

	initializeDatabase(db)
	gormDB = db
}

func initializeDatabase(db *gorm.DB) error {
	fmt.Println(">>> InitializeDatabase")
	// 테이블 생성
	if err := db.AutoMigrate(&member.MemberEntity{}, &site.SettingEntity{}); err != nil {
		return err
	}
	// siteadm 계정 만들기
	var signId string
	db.Raw("SELECT sign_id FROM members WHERE sign_id = ?", "siteadm").Scan(&signId)

	if len(signId) == 0 {
		db.Exec("INSERT INTO members(type, sign_id, name, password, created_at, updated_at) values(?, ?, ?, ?, datetime('now'), datetime('now'))",
			"site", "siteadm", "사이트 관리자", "$2a$04$7Ca1ybGc4yFkcBnzK1C0qevHy/LSD7PuBbPQTZEs6tiNM4hAxSYiG")
	}

	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Pre(echomiddleware.RemoveTrailingSlash())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowCredentials: true,
	}))

	e.Use(middlewares.GORMDb(gormDB))
	e.Use(middlewares.JwtToken())
	e.HideBanner = true

	controllers.AuthController{}.Init(e.Group("/api/auth"))
	controllers.SiteController{}.Init(e.Group("/api/site"))
	controllers.MemberController{}.Init(e.Group("/api/members"))

	color.Println(banner, color.Red("v"+Version), color.Blue(website))
	e.Start(":2016")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
