package sqlsugar

import (
	"database/sql"
	"fmt"
	"reflect"

	"strings"

	"github.com/pkg/errors"
)

type UpdateQuery struct {
	query          string
	setColumns     []string
	multipleTables bool
	err            error
}

var (
	NoSet                    = fmt.Errorf("Set not set")
	InappropriateSetAllUsage = fmt.Errorf("Inappropriate usage of SetAll method, can used only with single table queries")
)

func Update(table string) *UpdateQuery {
	return &UpdateQuery{
		query: fmt.Sprintf("UPDATE `%s`", table),
	}
}

func UpdateMultiple(tables []string) *UpdateQuery {
	quotedTables := []string{}
	for _, table := range tables {
		quotedTables = append(quotedTables, fmt.Sprintf("`%s`", table))
	}
	return &UpdateQuery{
		query:          fmt.Sprintf("UPDATE %s", strings.Join(quotedTables, ", ")),
		multipleTables: len(tables) > 1,
	}
}

func (q *UpdateQuery) Set(columns []string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(columns) == 0 {
		q.err = errors.WithStack(NoSet)
		return q
	}

	q.setColumns = []string{}
	setsExpr := ""
	for _, column := range columns {
		if len(setsExpr) > 0 {
			setsExpr += ", "
		}
		if idx := strings.Index(column, "."); idx != -1 {
			setsExpr += fmt.Sprintf("`%s`.`%s` = ?", column[:idx], column[idx+1:])
			q.setColumns = append(q.setColumns, column[idx+1:])
		} else {
			setsExpr += fmt.Sprintf("`%s` = ?", column)
			q.setColumns = append(q.setColumns, column)
		}
	}

	q.query = fmt.Sprintf("%s SET %s", q.query, setsExpr)
	return q
}

func (q *UpdateQuery) SetAll(i interface{}) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if q.multipleTables {
		q.err = errors.WithStack(InappropriateSetAllUsage)
		return q
	}

	columns := []string{}

	reflectedType := reflect.TypeOf(i).Elem()
	for i := 0; i < reflectedType.NumField(); i++ {
		fieldType := reflectedType.Field(i)

		column := fieldType.Tag.Get("column")
		if len(column) == 0 || column == "id" {
			continue
		}

		columns = append(columns, column)
	}

	return q.Set(columns)
}

func (q *UpdateQuery) Where(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s WHERE %s", q.query, condition)
	return q
}

func (q *UpdateQuery) OrderBy(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s ORDER BY %s", q.query, condition)
	return q
}

func (q *UpdateQuery) Limit(condition string) *UpdateQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s LIMIT %s", q.query, condition)
	return q
}

func (q *UpdateQuery) Exec(tx *sql.Tx, i interface{}, args ...interface{}) (sql.Result, error) {
	if len(q.setColumns) == 0 {
		return nil, errors.WithStack(NoSet)
	}

	fields := make([]interface{}, len(q.setColumns), len(q.setColumns))

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

	ex := fetchExecutor(tx)
	result, err := ex.Exec(q.query, sqlArgs...)
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

func (q *UpdateQuery) Error() error {
	return q.err
}
