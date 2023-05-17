package store

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/go-rest-template/internal/config"
	metricsMocks "github.com/zacscoding/go-rest-template/internal/metrics/mocks"
	"github.com/zacscoding/go-rest-template/internal/store/mocks"
	"github.com/zacscoding/go-rest-template/pkg/cache"
)

type CacheStoreSuite struct {
	suite.Suite
	conf         *config.Config
	cacher       cache.Cacher
	cacheCloseFn cache.CloseFn
	mpMock       *metricsMocks.Provider

	userStore     *userCacheStore
	userStoreMock *mocks.UserStore
}

func TestCacheStoreSuite(t *testing.T) {
	suite.Run(t, new(CacheStoreSuite))
}

func (s *CacheStoreSuite) BeforeTest(_, _ string) {
	conf, err := config.Load("", nil)
	s.NoError(err)

	s.conf = conf
	s.mpMock = &metricsMocks.Provider{}
	cacher, cacherCloseFn, err := cache.NewTestMemoryRedisCacher(s.T())
	s.NoError(err)
	s.cacher, s.cacheCloseFn = cacher, cacherCloseFn

	s.userStoreMock = &mocks.UserStore{}
	s.userStore = &userCacheStore{
		cacher:   s.cacher,
		mp:       s.mpMock,
		delegate: s.userStoreMock,
	}
}

func (s *CacheStoreSuite) AfterTest(_, _ string) {
	if s.cacheCloseFn != nil {
		s.cacheCloseFn()
	}
}
