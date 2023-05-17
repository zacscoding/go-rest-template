package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testWrapError(t *testing.T, db *gorm.DB) {
	name := "user1"
	assert.NoError(t, db.Create(&TestUser{Name: name}).Error)

	t.Run("ErrRecordNotFound", func(t *testing.T) {
		err := db.First(&TestUser{ID: 100}).Error

		assert.Error(t, err)
		assert.Equal(t, ErrRecordNotFound, WrapError(err))
	})

	t.Run("ErrKeyConflict", func(t *testing.T) {
		err := db.Create(&TestUser{Name: name}).Error

		assert.Error(t, err)
		assert.Equal(t, ErrKeyConflict, WrapError(err))
	})

	t.Run("ErrFKConstraint", func(t *testing.T) {
		err := db.Create(&TestCard{TestUserID: 200}).Error

		assert.Error(t, err)
		assert.Equal(t, ErrFKConstraint, WrapError(err))
	})
}
