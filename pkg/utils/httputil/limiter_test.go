package httputil

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLimiterRoundTripper(t *testing.T) {
	t.Run("NoRPS", func(t *testing.T) {
		got := NewLimiterRoundTripper(0, http.DefaultTransport)

		assert.Equal(t, http.DefaultTransport, got)
	})

	t.Run("WithRPS", func(t *testing.T) {
		got := NewLimiterRoundTripper(10, http.DefaultTransport)

		assert.NotEqual(t, http.DefaultTransport, got)
		limiter, ok := got.(*LimiterRoundTripper)
		assert.True(t, ok)
		assert.NotNil(t, limiter.limiter)
		assert.Equal(t, http.DefaultTransport, limiter.rt)
	})
}

func TestLimiterRoundTripper_RountTrip(t *testing.T) {
	rps := 5
	workers := rps * 2
	rt := NewLimiterRoundTripper(rps, http.DefaultTransport)
	cli := http.Client{Transport: rt}

	var (
		called = make(map[int64]int)
		mu     sync.Mutex
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		mu.Lock()
		called[now] = called[now] + 1
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cli.Get(ts.URL)
		}()
	}
	wg.Wait()

	total := 0
	for _, count := range called {
		assert.LessOrEqual(t, count, rps)
		total += count
	}
	assert.EqualValues(t, workers, total)
}
