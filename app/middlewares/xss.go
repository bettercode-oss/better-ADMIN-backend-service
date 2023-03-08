package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	ContentType                              = "Content-Type"
	ContentTypeApplicationJson               = "application/json"
	ContentTypeApplicationXWWWFormURLEncoded = "application/x-www-form-urlencoded"
)

var readFn = ioutil.ReadAll

func bytesToJsonMap(data []byte) (map[string]any, error) {
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, err
	}
	return jsonMap, nil
}
func jsonMapToBytes(j map[string]any) ([]byte, error) {
	bytesData, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return bytesData, nil
}
func sanitizeJsonMap(jsonMap map[string]any, p *bluemonday.Policy) map[string]any {
	sanitizedJsonMap := make(map[string]any)
	if jsonMap != nil {
		for key, val := range jsonMap {
			switch val.(type) {
			case string:
				str := fmt.Sprintf("%v", val)
				sanitizedJsonMap[key] = strings.TrimSpace(p.Sanitize(str))
			case float64:
				str := strconv.FormatFloat(val.(float64), 'g', 0, 64)
				numStr, err := strconv.ParseFloat(strings.TrimSpace(p.Sanitize(str)), 64)
				if err != nil {
					log.Error(err.Error())
					break
				}
				sanitizedJsonMap[key] = numStr
			case map[string]any:
				sanitizedJsonMap[key] = sanitizeJsonMap(val.(map[string]any), p)
			default:
				sanitizedJsonMap[key] = strings.TrimSpace(p.Sanitize(fmt.Sprintf("%v", val)))
			}
		}
	}
	return sanitizedJsonMap
}
func XssSanitizer(httpMethods ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		contentType := req.Header.Get(ContentType)
		if slices.Contains(httpMethods, req.Method) {
			if contentType == ContentTypeApplicationJson {
				if req.Body != nil {
					bodyBytes, err := readFn(req.Body)
					if err != nil {
						log.Errorf("%+v", err.Error())
						c.Next()
						return
					}
					jsonMap, err := bytesToJsonMap(bodyBytes)
					if err != nil {
						log.Errorf("%+v", err.Error())
						req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
						c.Next()
						return
					}
					jsonBytes, err := jsonMapToBytes(sanitizeJsonMap(jsonMap, bluemonday.UGCPolicy()))
					if err != nil {
						log.Errorf("%+v", err.Error())
						req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
						c.Next()
						return
					}
					req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))
				}
			}
		}
	}
}
