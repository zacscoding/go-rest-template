package store

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/metrics"
	"github.com/zacscoding/go-rest-template/pkg/database"
	"gorm.io/gorm"
)

const migrationDir = "../../migrations"

type StoreSuite struct {
	suite.Suite
	conf    *config.Config
	dsn     string
	db      *gorm.DB
	closeFn database.CloseFunc

	userStore UserStore
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

func (s *StoreSuite) SetupSuite() {
	conf, err := config.Load("", nil)
	s.NoError(err)

	s.conf = conf
	mp := metrics.NewProvider(s.conf)
	s.dsn, s.db, s.closeFn = database.NewTestMysqlDB(s.T(), "")
	s.userStore, _ = NewUserStore(nil, s.db, nil, mp)
}

func (s *StoreSuite) BeforeTest(_, _ string) {
	s.NoError(database.MigrateMysqlDB(s.dsn, migrationDir, true))
}

func (s *StoreSuite) AfterTest(_, _ string) {
	s.NoError(database.MigrateMysqlDB(s.dsn, migrationDir, false))
}

func (s *StoreSuite) TearDownSuite() {
	if s.closeFn != nil {
		s.NoError(s.closeFn())
	}
}
