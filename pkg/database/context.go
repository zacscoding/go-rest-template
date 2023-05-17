package database

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type contextKey = string

const dbKey = contextKey("db")

// FromContext returns the *gorm.DB stored in the context if exists, otherwise returns defaultDB.
func FromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if ctx == nil {
		return defaultDB
	}
	if db, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		return db
	}
	return defaultDB
}

// WithContext creates a new context with the provided db attached.
func WithContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}
