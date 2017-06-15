package sqlsugar

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type DeleteQuery struct {
	query string
	err   error
}

func Delete(table string) *DeleteQuery {
	return &DeleteQuery{
		query: fmt.Sprintf("DELETE FROM `%s`", table),
	}
}

func (q *DeleteQuery) Where(condition string) *DeleteQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &DeleteQuery{
		query: fmt.Sprintf("%s WHERE %s", q.query, condition),
	}
}

func (q *DeleteQuery) OrderBy(condition string) *DeleteQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s ORDER BY %s", q.query, condition)
	return q
}

func (q *DeleteQuery) Limit(condition string) *DeleteQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s LIMIT %s", q.query, condition)
	return q
}

func (q *DeleteQuery) Exec(tx *sql.Tx, args ...interface{}) (sql.Result, error) {
	ex := fetchExecutor(tx)
	result, err := ex.Exec(q.query, args...)
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (q *DeleteQuery) Error() error {
	return q.err
}
