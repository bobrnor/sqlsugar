package sqlsugar

import (
	"database/sql"
	"time"

	"context"

	"github.com/pkg/errors"
)

type TxRollbackFunc func(err error)

type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

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
		err := database.Close()
		return errors.Wrap(err, "Can't close database")
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

func Begin() (*sql.Tx, error) {
	tx, err := database.Begin()
	return tx, errors.WithStack(err)
}

func RollbackOnRecover(tx *sql.Tx, fn TxRollbackFunc) {
	if err := recover(); err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr == nil && fn != nil {
			fn(errors.Errorf("%+v", err))
		}
	}
}

func fetchExecutor(tx *sql.Tx) executor {
	if tx == nil {
		return database
	}
	return tx
}
