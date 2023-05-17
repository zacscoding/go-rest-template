package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zacscoding/go-rest-template/internal/handler/apierr"
)

func TestHandleResponse(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name string
		data interface{}
		err  error
		// expected
		status int
		body   string
	}{
		{
			name: "Success",
			data: map[string]interface{}{
				"key":    "key1",
				"pubkey": "mypubkey",
			},
			status: http.StatusOK,
			body: `{
			  "key": "key1",
			  "pubkey": "mypubkey"
			}`,
		},
		{
			name:   "Status Error",
			err:    apierr.ErrInvalidRequest.WithStatusCode(http.StatusUnauthorized),
			status: http.StatusUnauthorized,
			body: `{
			  "code": "InvalidRequest",
			  "message": "Request form is not valid."
			}`,
		},
		{
			name:   "Status Error With Code",
			err:    apierr.ErrInvalidRequest.WithCode("CustomCode"),
			status: http.StatusBadRequest,
			body: `{
			  "code": "CustomCode",
			  "message": "Request form is not valid."
			}`,
		},
		{
			name:   "Unknown Error",
			err:    errors.New("internal server error"),
			status: http.StatusInternalServerError,
			body: `{
			  "code": "InternalServerError",
			  "message": "There was an error. Please try again later."
			}`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			r.GET("/foo", func(ctx *gin.Context) {
				HandleResponse(ctx, tc.data, tc.err)
			})

			res := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "http://localhost/foo", nil)

			r.ServeHTTP(res, req)

			assert.Equal(t, tc.status, res.Code)
			assert.JSONEq(t, tc.body, res.Body.String())
		})
	}
}
