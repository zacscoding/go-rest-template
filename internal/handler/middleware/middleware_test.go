package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		requestID string
	}{
		{name: "Exist RequestID", requestID: "request1"},
		{name: "Empty RequestID", requestID: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			srv := setupRouter(func(c *gin.Engine) {
				c.Use(RequestIDMiddleware())
			})

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "http://localhost/foo", nil)
			if tc.requestID != "" {
				req.Header.Set(XRequestIdKey, tc.requestID)
			}

			srv.ServeHTTP(res, req)

			if tc.requestID == "" {
				assert.NotEmpty(t, res.Header().Get(XRequestIdKey))
			} else {
				assert.Equal(t, tc.requestID, res.Header().Get(XRequestIdKey))
			}
		})
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	timeout := time.Millisecond * 50
	srv := setupRouterWithHandler(func(c *gin.Engine) {
		c.Use(TimeoutMiddleware(timeout))
	}, func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		assert.True(t, ok)
		assert.LessOrEqual(t, time.Since(deadline), timeout)
		time.Sleep(500 * time.Millisecond)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost/foo", nil)

	srv.ServeHTTP(res, req)

	assert.Equal(t, http.StatusGatewayTimeout, res.Code)
}

func setupRouter(middlewareFunc func(c *gin.Engine)) *gin.Engine {
	return setupRouterWithHandler(middlewareFunc, func(c *gin.Context) {
		c.JSON(200, "bar")
	})
}

func setupRouterWithHandler(middlewareFunc func(c *gin.Engine), handler func(c *gin.Context)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	middlewareFunc(r)
	r.GET("/foo", handler)
	return r
}
