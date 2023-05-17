package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/zacscoding/go-rest-template/pkg/logging"
	"go.uber.org/zap/zapcore"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func openMysqlDB(conf *Config) (*gorm.DB, error) {
	var (
		db     *gorm.DB
		err    error
		logger = NewLogger(time.Second, true, zapcore.Level(conf.LoggingLevel), conf.LoggingPrefix)
	)

	for i := 0; i < 20; i++ {
		db, err = gorm.Open(gmysql.Open(conf.DataSourceName), &gorm.Config{Logger: logger})
		if err == nil {
			break
		}
		logging.DefaultLogger().Warnf("failed to open database: %v", err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rawDB.SetMaxOpenConns(conf.Pool.MaxOpen)
	rawDB.SetMaxIdleConns(conf.Pool.MaxIdle)
	rawDB.SetConnMaxLifetime(conf.Pool.MaxLifeTime)

	var replicas []gorm.Dialector
	for _, dsn := range conf.Replica.DataSourceNames {
		replicas = append(replicas, gmysql.Open(dsn))
	}
	if len(replicas) != 0 {
		if err := db.Use(
			dbresolver.Register(dbresolver.Config{Replicas: replicas}).
				SetMaxOpenConns(conf.Replica.Pool.MaxOpen).
				SetMaxIdleConns(conf.Replica.Pool.MaxIdle).
				SetConnMaxLifetime(conf.Replica.Pool.MaxLifeTime),
		); err != nil {
			return nil, fmt.Errorf("register replica resolvers: %v", err)
		}
	}

	if conf.Migrate.Enabled {
		err := MigrateMysqlDB(conf.DataSourceName, conf.Migrate.Dir, true)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

// MigrateMysqlDB migrates database from the given dsn data source name and migration directories.
func MigrateMysqlDB(dsn, dir string, isUp bool) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed create connect database: %w", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to mysql instance: %w", err)
	}
	if dir == "" {
		return ErrEmptyDir
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", dir),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to new database instance: %w", err)
	}
	var migrateErr error
	if isUp {
		migrateErr = m.Up()
	} else {
		migrateErr = m.Down()
	}
	if migrateErr != nil && !errors.Is(migrateErr, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", migrateErr)
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed close db: %w", dbErr)
	}
	return nil
}

// NewTestMysqlDB starts mysql docker container with given version tag and returns dsn, gorm.DB, CloseFunc to clean up.
// If provide empty tag, will use "5.7".
func NewTestMysqlDB(tb testing.TB, tag string) (string, *gorm.DB, CloseFunc) {
	tb.Helper()
	var db *sql.DB
	if tag == "" {
		tag = "8.0.30"
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Failed to connect to docker: %v", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", tag, []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		tb.Fatalf("Failed to not start resource: %v", err)
	}
	err = resource.Expire(60 * 5)
	if err != nil {
		tb.Fatalf("Failed to expire resource: %v", err)
	}

	dsnFmt := "root:secret@(localhost:%s)/mysql?charset=utf8&parseTime=true&multiStatements=true"
	dsn := fmt.Sprintf(dsnFmt, resource.GetPort("3306/tcp"))
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Failed to connect to docker: %v", err)
	}

	gdb, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to create a new gorm.DB: %s", err)
	}

	closeFn := func() error {
		_ = db.Close()
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Failed to purge resource: %s", err)
			return err
		}
		return nil
	}
	return dsn, gdb, closeFn
}
