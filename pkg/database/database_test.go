package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestUser struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"uniqueIndex;size:50"`
	TestCards []TestCard
}

type TestCard struct {
	ID         uint `gorm:"primarykey"`
	TestUserID uint
}

func testRunInTx(t *testing.T, db *gorm.DB) {
	name := "user1"
	assert.NoError(t, db.Create(&TestUser{Name: name}).Error)

	err := RunInTx(context.TODO(), db, nil, func(txDb *gorm.DB) error {
		if err := txDb.Create(&TestUser{Name: name + "_1"}).Error; err != nil {
			return err
		}
		if err := txDb.Create(&TestUser{Name: name + "_2"}).Error; err != nil {
			return err
		}
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, db.Where("name = ?", name+"_1").First(new(TestUser)).Error)
	assert.NoError(t, db.Where("name = ?", name+"_2").First(new(TestUser)).Error)
}

func testRunInTxRollback(t *testing.T, db *gorm.DB) {
	name := "user1"
	assert.NoError(t, db.Create(&TestUser{Name: name}).Error)
	firstSuccess := false

	err := RunInTx(context.TODO(), db, nil, func(txDb *gorm.DB) error {
		if err := txDb.Create(&TestUser{Name: name + "_1"}).Error; err != nil {
			return err
		}
		firstSuccess = true
		if err := txDb.Create(&TestUser{Name: name}).Error; err != nil {
			return err
		}
		return nil
	})

	assert.Error(t, err)
	assert.True(t, firstSuccess)
	assert.Equal(t, ErrRecordNotFound, WrapError(db.Where("name = ?", name+"_1").First(&TestUser{}).Error))
}
