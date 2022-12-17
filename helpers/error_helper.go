package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	errorHelperOnce     sync.Once
	errorHelperInstance *errorHelper
)

func ErrorHelper() *errorHelper {
	errorHelperOnce.Do(func() {
		errorHelperInstance = &errorHelper{}
	})

	return errorHelperInstance
}

type errorHelper struct {
}

func (errorHelper) InternalServerError(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Status(http.StatusInternalServerError)
}
