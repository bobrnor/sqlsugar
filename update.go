package sqlsugar

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type UpdateQuery struct {
	query      string
	setColumns []string
	err        error
}

var (
	NoSet = errors.New("Set not set")
)

func Update(table string) *UpdateQuery {
	return &UpdateQuery{
		query: fmt.Sprintf("UPDATE `%s`", table),
	}
}

func (q *UpdateQuery) Set(columns []string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(columns) == 0 {
		return &UpdateQuery{
			query: "",
			err:   errors.WithStack(NoSet),
		}
	}

	setsExpr := ""
	for _, column := range columns {
		if len(setsExpr) > 0 {
			setsExpr += ", "
		}

		setsExpr += fmt.Sprintf("`%s` = ?", column)
	}

	q.setColumns = columns

	return &UpdateQuery{
		query: fmt.Sprintf("%s SET %s", q.query, setsExpr),
	}
}

func (q *UpdateQuery) Where(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &UpdateQuery{
		query: fmt.Sprintf("%s WHERE %s", q.query, condition),
	}
}

func (q *UpdateQuery) OrderBy(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &UpdateQuery{
		query: fmt.Sprintf("%s ORDER BY %s", q.query, condition),
	}
}

func (q *UpdateQuery) Limit(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &UpdateQuery{
		query: fmt.Sprintf("%s LIMIT %s", q.query, condition),
	}
}

func (q *UpdateQuery) Exec(tx *sql.Tx, i interface{}, args ...interface{}) (sql.Result, error) {
	if len(q.setColumns) == 0 {
		return nil, errors.WithStack(NoSet)
	}

	fields := make([]interface{}, len(q.setColumns), 0)

	reflectedType := reflect.TypeOf(i).Elem()
	reflectedValue := reflect.ValueOf(i).Elem()
	for i := 0; i < reflectedValue.NumField(); i++ {
		fieldType := reflectedType.Field(i)
		fieldValue := reflectedValue.Field(i)

		column := fieldType.Tag.Get("column")
		if len(column) == 0 {
			continue
		}

		if idx := index(q.setColumns, column); idx != -1 {
			fields[idx] = fieldValue.Interface()
		}
	}

	sqlArgs := append(fields, args...)

	result, err := database.Exec(q.query, sqlArgs...)
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func index(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
