package middlewares

import (
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestAuthorizer_권한이_있는_API_접근할_수_있다(t *testing.T) {
	// given
	router := gin.Default()
	opaRego, err := rego.New(rego.Query("data.rest.allowed"), rego.Load([]string{
		"../../authorization/rest/policy.rego", "../../authorization/rest/data.json",
	}, nil)).PrepareForEval(context.TODO())

	if err != nil {
		t.Error(err)
	}

	router.Use(JwtToken())
	router.Use(RestAuthorizer(&opaRego))
	router.GET("/api/access-control/permissions", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions", nil)
	req = req.WithContext(helpers.ContextHelper().SetUserClaim(req.Context(), &security.UserClaim{
		Id:          1,
		Permissions: []string{"access-control-permission.read"},
	}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRestAuthorizer_권한이_없는_API에_접근할_수_없다(t *testing.T) {
	// given
	router := gin.Default()
	opaRego, err := rego.New(rego.Query("data.rest.allowed"), rego.Load([]string{
		"../../authorization/rest/policy.rego", "../../authorization/rest/data.json",
	}, nil)).PrepareForEval(context.TODO())

	if err != nil {
		t.Error(err)
	}

	router.Use(JwtToken())
	router.Use(RestAuthorizer(&opaRego))
	router.GET("/api/access-control/permissions", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions", nil)
	req = req.WithContext(helpers.ContextHelper().SetUserClaim(req.Context(), &security.UserClaim{
		Id:          1,
		Permissions: []string{},
	}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
