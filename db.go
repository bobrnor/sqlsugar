package sorm

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

type TxRollbackFunc func(err error)

var (
	database *sql.DB
)

func Open(driverName, dataSourceName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return errors.WithStack(err)
	}

	database = db
	return nil
}

func Close() error {
	if database != nil {
		return database.Close()
	}
	return nil
}

func SetMaxOpenConns(n int) {
	if database != nil {
		database.SetMaxOpenConns(n)
	}
}

func SetMaxIdleConns(n int) {
	if database != nil {
		database.SetMaxIdleConns(n)
	}
}

func SetConnMaxLifetime(d time.Duration) {
	if database != nil {
		database.SetConnMaxLifetime(d)
	}
}

func RollbackOnRecover(tx *sql.Tx, fn TxRollbackFunc) {
	if err := recover(); err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr == nil && err.(error) != nil && fn != nil {
			fn(errors.WithStack(err.(error)))
		}
	}
}
