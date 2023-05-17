package cache

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RedisCacheSuite struct {
	suite.Suite
	cacher  Cacher
	closeFn CloseFn
}

func TestRedisCache(t *testing.T) {
	suite.Run(t, new(RedisCacheSuite))
}

func (s *RedisCacheSuite) SetupSuite() {
	var err error
	// TODO: create a new cacher depends on testing tags
	s.cacher, s.closeFn, err = NewTestMemoryRedisCacher(s.T())
	// s.cacher, s.closeFn, err = NewTestRedisCacher(s.T())
	// s.cacher, s.closeFn, err = NewTestClusterRedisCacher(s.T())
	s.NoError(err)
}

func (s *RedisCacheSuite) TearDownSuite() {
	if s.closeFn != nil {
		s.closeFn()
	}
}

func (s *RedisCacheSuite) TestFetch() {
	testFetch(s.T(), s.cacher)
}

func (s *RedisCacheSuite) TestGet() {
	testGet(s.T(), s.cacher)
}

func (s *RedisCacheSuite) TestSet() {
	testSet(s.T(), s.cacher)
}

func (s *RedisCacheSuite) TestExists() {
	testExists(s.T(), s.cacher)
}

func (s *RedisCacheSuite) TestDelete() {
	testDelete(s.T(), s.cacher)
}
