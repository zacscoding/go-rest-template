package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/zacscoding/go-rest-template/internal/handler/apierr"
	"github.com/zacscoding/go-rest-template/internal/handler/middleware"
)

func Wrap(fn func(gctx *gin.Context) (data interface{}, err error)) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		data, err := fn(gctx)
		HandleResponse(gctx, data, err)
	}
}

func HandleResponse(gctx *gin.Context, data interface{}, err error) {
	if data == nil && err == nil {
		gctx.Status(http.StatusOK)
		return
	}
	if err == nil {
		gctx.JSON(http.StatusOK, data)
		return
	}

	errorResp := apierr.ErrInternalServerError
	if e, ok := err.(*apierr.Error); ok {
		errorResp = e
	}
	errorResp.RequestID = gctx.Writer.Header().Get(middleware.XRequestIdKey)
	gctx.AbortWithStatusJSON(errorResp.StatusCode, errorResp)
}
