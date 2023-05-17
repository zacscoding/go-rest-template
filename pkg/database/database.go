package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUnsupportedDriver = errors.New("unsupported database driver")
	ErrEmptyDir          = errors.New("empty directory")
)

type CloseFunc func() error

// Config represents database configs.
type Config struct {
	Driver         string `json:"driver" yaml:"driver"`
	DataSourceName string `json:"data-source-name" yaml:"data-source-name"`
	LoggingLevel   int    `json:"logging-level" yaml:"logging-level"`
	LoggingPrefix  string `json:"logging-prefix" yaml:"logging-prefix"`
	BatchSize      int    `json:"batch-size" yaml:"batch-size"`
	Migrate        struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Dir     string `json:"dir" yaml:"dir"`
	} `json:"migrate" yaml:"migrate"`
	Pool struct {
		MaxOpen     int           `json:"max-open" yaml:"max-open"`
		MaxIdle     int           `json:"max-idle" yaml:"max-idle"`
		MaxLifeTime time.Duration `json:"max-lifetime" yaml:"max-lifetime"`
	} `json:"pool" yaml:"pool"`
	Replica struct {
		DataSourceNames []string `json:"data-source-names" yaml:"data-source-names"`
		Pool            struct {
			MaxOpen     int           `json:"max-open" yaml:"max-open"`
			MaxIdle     int           `json:"max-idle" yaml:"max-idle"`
			MaxLifeTime time.Duration `json:"max-lifetime" yaml:"max-lifetime"`
		} `json:"pool" yaml:"pool"`
	} `json:"replica" yaml:"replica"`
}

// Open returns a new gorm.DB for given conf Config.
func Open(conf *Config) (*gorm.DB, error) {
	switch conf.Driver {
	case "mysql":
		return openMysqlDB(conf)
	default:
		return nil, ErrUnsupportedDriver
	}
}

var DefaulTxOptions = &sql.TxOptions{
	Isolation: sql.LevelDefault,
	ReadOnly:  false,
}

// RunInTx begin transaction from given database and execute f.
func RunInTx(ctx context.Context, db *gorm.DB, opts *sql.TxOptions, f func(txdb *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin(opts)
	if tx.Error != nil {
		return fmt.Errorf("start tx: %v", tx.Error)
	}

	if err := f(tx); err != nil {
		if err1 := tx.Rollback().Error; err1 != nil {
			return fmt.Errorf("rollback tx: %v (original error: %v)", err1, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}
	return nil
}
