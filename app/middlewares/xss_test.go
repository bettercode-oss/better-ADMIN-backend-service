package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func mirrorEntityBody(ctx *gin.Context) {
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	ctx.JSON(http.StatusOK, body)
}
func TestXssSanitizer_json_string_과_number_필드_타입_살균_성공(t *testing.T) {
	// given
	router := gin.Default()
	router.POST("/", XssSanitizer(http.MethodPost, http.MethodPut), mirrorEntityBody)
	requestBody := `{
		"id": 1,
	    "data": "hello <script>alert('xss attack');</script>",
        "html" :"<a onblur='alert(secret)' href='http://www.google.com'>Google</a>"
	}`
	// when
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(ContentType, ContentTypeApplicationJson)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id":   float64(1),
		"data": "hello",
		"html": `<a href="http://www.google.com" rel="nofollow">Google</a>`,
	}
	assert.Equal(t, expected, actual)
}
func TestXssSanitizer_중첩된_json_살균_성공(t *testing.T) {
	// given
	router := gin.Default()
	router.POST("/", XssSanitizer(http.MethodPost, http.MethodPut), mirrorEntityBody)
	requestBody := `{
		"id":   2,
        "htmlJson" : {
       		 "attack" :"<a onblur='alert(secret)' href='http://www.google.com'>Google</a>"
		}
	}`
	// when
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(ContentType, ContentTypeApplicationJson)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id": float64(2),
		"htmlJson": map[string]any{
			"attack": `<a href="http://www.google.com" rel="nofollow">Google</a>`,
		},
	}
	assert.Equal(t, expected, actual)
}
func TestXssSanitizer_requestBody_읽기_실패시_살균_안하고_다음_코드_실행(t *testing.T) {
	// given
	oldReadFn := readFn
	defer func() {
		readFn = oldReadFn
	}()
	readFn = func(r io.Reader) ([]byte, error) {
		return nil, errors.New(`err`)
	}
	router := gin.Default()
	router.POST("/", XssSanitizer(http.MethodPost, http.MethodPut), mirrorEntityBody)
	requestBody := `{
		"id": 2,
	    "data": "hello <script>alert('xss attack');</script>"
	}`
	// when
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(ContentType, ContentTypeApplicationJson)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id":   float64(2),
		"data": "hello <script>alert('xss attack');</script>",
	}
	assert.Equal(t, expected, actual)
}
func TestXssSanitizer_ContentTypeApplicationXWWWFormURLEncoded_살균_안함(t *testing.T) {
	// given
	router := gin.Default()
	router.POST("/", XssSanitizer(http.MethodPost, http.MethodPut), func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		password := ctx.PostForm("password")
		ctx.JSON(http.StatusOK, map[string]any{
			"id":       id,
			"password": password,
		})
	})

	formData := url.Values{
		"id":       {"test"},
		"password": {"3eq*81<>H2<Y9>t9"},
	}
	// when
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Add(ContentType, ContentTypeApplicationXWWWFormURLEncoded)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id":       "test",
		"password": "3eq*81<>H2<Y9>t9",
	}
	assert.Equal(t, expected, actual)
}

// TODO 특정 필드 살균 안하도록 수정
func TestXssSanitizer_json_특정_필드_살균_안함(t *testing.T) {
	// given
	router := gin.Default()
	router.POST("/", XssSanitizer(http.MethodPost, http.MethodPut), mirrorEntityBody)
	requestBody := `{
		"id": "test",
	    "password": "3eq*81<>H2<Y9>t9"
	}`
	// when
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Add(ContentType, ContentTypeApplicationJson)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id":       "test",
		"password": "3eq*81<>H2<Y9>t9",
	}
	assert.Equal(t, expected, actual)
}
func TestXssSanitizer_http_method_Patch_살균_안함(t *testing.T) {
	// given
	router := gin.Default()
	router.PATCH("/", XssSanitizer(http.MethodPost, http.MethodPut), mirrorEntityBody)
	requestBody := `{
		"id":   2,
        "htmlJson" : {
       		 "attack" :"<a onblur='alert(secret)' href='http://www.google.com'>Google</a>"
		}
	}`
	// when
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestBody))
	req.Header.Set(ContentType, ContentTypeApplicationJson)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// then
	assert.Equal(t, http.StatusOK, w.Code)
	actual, err := bytesToJsonMap(w.Body.Bytes())
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
	expected := map[string]any{
		"id": float64(2),
		"htmlJson": map[string]any{
			"attack": "<a onblur='alert(secret)' href='http://www.google.com'>Google</a>",
		},
	}
	assert.Equal(t, expected, actual)
}
func Test_jsonMapToBytes_실패(t *testing.T) {
	jsonMap := map[string]any{
		"data": make(chan int),
	}
	_, err := jsonMapToBytes(jsonMap)
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
}
func Test_bytesToJsonMap_실패(t *testing.T) {
	jsonMap := `{
		"id": 1
	    "data": "hello"
	}`
	_, err := bytesToJsonMap([]byte(jsonMap))
	if err != nil {
		log.Errorln(err.Error())
		t.Failed()
	}
}
