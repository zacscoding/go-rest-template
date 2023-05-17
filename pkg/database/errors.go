package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	// ErrRecordNotFound is an error if not exists record.
	ErrRecordNotFound = errors.New("record not found")
	// ErrKeyConflict is an error if duplicate field.
	ErrKeyConflict = errors.New("conflict key")
	// ErrFKConstraint is an error if foreign key constraint failed.
	ErrFKConstraint = errors.New("a foreign key constraint fails")
)

// WrapError wraps given error to database error.
// ErrRecordNotFound is returned if err is a gorm.ErrRecordNotFound
// If conflict key error, then ErrKeyConflict will be returned.
func WrapError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return ErrRecordNotFound
	}
	if e, ok := err.(*mysql.MySQLError); ok {
		return wrapMySQLError(e)
	}
	return err
}

func wrapMySQLError(err *mysql.MySQLError) error {
	switch err.Number {
	case 1062:
		return ErrKeyConflict
	case 1452:
		return ErrFKConstraint
	default:
		return err
	}
}
