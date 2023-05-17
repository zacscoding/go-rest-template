package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func TestFromContext(t *testing.T) {
	defaultDB := gorm.DB{}
	t.Run("no db in context", func(t *testing.T) {
		db1 := FromContext(context.Background(), &defaultDB)

		assert.Equal(t, &defaultDB, db1)
	})

	t.Run("db in context", func(t *testing.T) {
		newDB := gorm.DB{}

		ctx := WithContext(context.Background(), &newDB)
		db1 := FromContext(ctx, &defaultDB)

		assert.Equal(t, &newDB, db1)
	})
}
