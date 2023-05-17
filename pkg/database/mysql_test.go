package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var migrationTables = []string{
	"mysql_migration_users",
	"mysql_migration_users2",
	"schema_migrations",
}

type MysqlSuite struct {
	suite.Suite
	dsn     string
	db      *gorm.DB
	closeFn CloseFunc
}

func TestMysqlSuite(t *testing.T) {
	suite.Run(t, new(MysqlSuite))
}

func (s *MysqlSuite) SetupSuite() {
	s.dsn, s.db, s.closeFn = NewTestMysqlDB(s.T(), "")
}

func (s *MysqlSuite) SetupTest() {
	s.NoError(s.db.Migrator().AutoMigrate(new(TestUser), new(TestCard)))
}

func (s *MysqlSuite) TearDownTest() {
	s.NoError(s.db.Migrator().DropTable(new(TestUser), new(TestCard)))
	for _, table := range migrationTables {
		s.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	}
}

func (s *MysqlSuite) TearDownSuite() {
	if s.closeFn != nil {
		err := s.closeFn()
		s.NoError(err)
	}
}

func (s *MysqlSuite) TestOpenMysqlDB() {
	var conf Config
	conf.Driver = "mysql"
	conf.DataSourceName = s.dsn
	conf.Migrate.Enabled = true
	conf.Migrate.Dir = "./migrations/mysql"
	conf.Pool.MaxOpen = 20
	conf.Pool.MaxIdle = 15
	conf.Pool.MaxLifeTime = time.Minute

	db, err := openMysqlDB(&conf)

	s.NoError(err)
	expectedTables := []string{
		"mysql_migration_users",
		"mysql_migration_users2",
		"schema_migrations",
	}
	tables, err := db.Migrator().GetTables()
	s.NoError(err)
	for _, e := range expectedTables {
		s.Contains(tables, e)
	}

	sqlDB, err := db.DB()
	s.NoError(err)
	stats := sqlDB.Stats()
	s.EqualValues(conf.Pool.MaxOpen, stats.MaxOpenConnections)
}

func (s *MysqlSuite) TestRunInTx() {
	testRunInTx(s.T(), s.db)
}

func (s *MysqlSuite) TestRunInTx_Rollback() {
	testRunInTxRollback(s.T(), s.db)
}

func (s *MysqlSuite) TestWrapError() {
	testWrapError(s.T(), s.db)
}

func (s *MysqlSuite) TestMigrateMysqlDB() {
	err := MigrateMysqlDB(s.dsn, "./migrations/mysql", true)

	s.NoError(err)
	tables, err := s.db.Migrator().GetTables()
	s.NoError(err)
	for _, table := range migrationTables {
		s.Contains(tables, table)
	}
}
