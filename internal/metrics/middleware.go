package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewMiddleware(mp Provider, skipPaths ...string) gin.HandlerFunc {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, pth := range skipPaths {
		skip[pth] = struct{}{}
	}

	return func(gctx *gin.Context) {
		// skip record metrics
		if _, ok := skip[gctx.FullPath()]; ok {
			gctx.Next()
			return
		}

		start := time.Now()
		gctx.Next()
		elapsd := time.Since(start)
		var (
			code   = gctx.Writer.Status()
			method = gctx.Request.Method
			path   = gctx.FullPath()
		)
		mp.RecordApiCount(code, method, path)
		mp.RecordApiLatency(code, method, path, elapsd)
	}
}
