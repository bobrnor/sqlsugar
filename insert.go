package sqlsugar

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type InsertQuery struct {
	query string
	t     reflect.Type

	tableSet bool

	err error
}

var (
	NoTable = errors.New("No table")
)

func Insert(i interface{}) *InsertQuery {
	reflectedType := reflect.TypeOf(i).Elem()

	columns := []string{}
	values := []string{}
	for i := 0; i < reflectedType.NumField(); i++ {
		field := reflectedType.Field(i)

		column := field.Tag.Get("column")
		if len(column) == 0 || column == "id" {
			continue
		}

		quotedColumn := fmt.Sprintf("`%s`", column)
		columns = append(columns, quotedColumn)
		values = append(values, "?")
	}

	columnsExpr := strings.Join(columns, ", ")
	valuesExpr := strings.Join(values, ", ")
	return &InsertQuery{
		query: fmt.Sprintf("(%s) VALUES (%s)", columnsExpr, valuesExpr),
		t:     reflectedType,
	}
}

func (q *InsertQuery) Into(table string) *InsertQuery {
	if q.err != nil {
		return q
	}

	if len(table) == 0 {
		q.err = errors.WithStack(NoTable)
		return q
	}

	q.tableSet = true
	q.query = fmt.Sprintf("INSERT INTO `%s` %s", table, q.query)
	return q
}

func (q *InsertQuery) Exec(tx *sql.Tx, i interface{}) (sql.Result, error) {
	if !q.tableSet {
		return nil, errors.WithStack(NoTable)
	}

	args := []interface{}{}

	reflectedType := reflect.TypeOf(i).Elem()
	reflectedValue := reflect.ValueOf(i).Elem()
	for i := 0; i < reflectedValue.NumField(); i++ {
		fieldType := reflectedType.Field(i)
		fieldValue := reflectedValue.Field(i)

		column := fieldType.Tag.Get("column")
		if len(column) == 0 || column == "id" {
			continue
		}

		args = append(args, fieldValue.Interface())
	}

	ex := fetchExecutor(tx)
	result, err := ex.Exec(q.query, args...)
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (q *InsertQuery) Error() error {
	return q.err
}
