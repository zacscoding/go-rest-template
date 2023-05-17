package httputil

import (
	"net/http"

	"go.uber.org/ratelimit"
)

type LimiterRoundTripper struct {
	limiter ratelimit.Limiter
	rt      http.RoundTripper
}

// NewLimiterRoundTripper creates a new LimiterRoundTripper for given rps and http.RoundTripper.
// If rps is less than or equals to 0, return given rt http.RoundTripper.
func NewLimiterRoundTripper(rps int, rt http.RoundTripper) http.RoundTripper {
	if rps <= 0 {
		return rt
	}
	return &LimiterRoundTripper{limiter: ratelimit.New(rps), rt: rt}
}

func (lrt *LimiterRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	lrt.limiter.Take()
	return lrt.rt.RoundTrip(req)
}
