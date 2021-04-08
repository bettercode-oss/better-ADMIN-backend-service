package controllers

import (
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/middlewares"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/testfixtures.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}


func TestAuthController_AuthWithSignIdPassword(t *testing.T) {
	// given
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(&member.MemberEntity{})

	sqlDB, _ := gormDB.DB()

	fixtures, err := testfixtures.NewFolder(sqlDB, &testfixtures.SQLite{}, "../testdata/db_fixtures")
	if err != nil {
		panic(err)
	}

	testfixtures.SkipDatabaseNameCheck(true)

	if err := fixtures.Load(); err != nil {
		panic(err)
	}

	db := middlewares.GORMDb(gormDB)

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return db(handlerFunc)(c)
	}

	requestBody := `{
		"id": "siteadm",
		"password": "123456"
	}`

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// when
	handleWithFilter(AuthController{}.AuthWithSignIdPassword, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.(map[string]interface{})["jwtToken"])
}
