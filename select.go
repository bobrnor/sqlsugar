package sqlsugar

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type SelectQuery struct {
	query string
	t     reflect.Type

	tableSet bool

	err error
}

type row interface {
	Scan(dest ...interface{}) error
}

type rows interface {
	row
	Next() bool
	Err() error
}

var (
	NoTablesUsed = errors.New("No tables used")
)

func Select(i interface{}) *SelectQuery {
	selectExpr := selectExpression(i)
	return &SelectQuery{
		query: fmt.Sprintf("SELECT %s", selectExpr),
		t:     reflect.TypeOf(i).Elem(),
	}
}

func (q *SelectQuery) From(tables []string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(tables) == 0 {
		q.err = errors.WithStack(NoTablesUsed)
		return q
	}

	quotedTables := []string{}
	for _, table := range tables {
		quotedTables = append(quotedTables, fmt.Sprintf("`%s`", table))
	}

	q.tableSet = true
	q.query = fmt.Sprintf("%s FROM %s", q.query, strings.Join(quotedTables, ", "))
	return q
}

func (q *SelectQuery) Where(condition string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s WHERE %s", q.query, condition)
	return q
}

func (q *SelectQuery) GroupBy(condition string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s GROUP BY %s", q.query, condition)
	return q
}

func (q *SelectQuery) Having(condition string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s HAVING %s", q.query, condition)
	return q
}

func (q *SelectQuery) OrderBy(condition string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s ORDER BY %s", q.query, condition)
	return q
}

func (q *SelectQuery) Limit(condition string) *SelectQuery {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	q.query = fmt.Sprintf("%s LIMIT %s", q.query, condition)
	return q
}

func (q *SelectQuery) ForUpdate() *SelectQuery {
	if q.err != nil {
		return q
	}

	q.query = fmt.Sprintf("%s FOR UPDATE", q.query)
	return q
}

func selectExpression(i interface{}) string {
	reflectedType := reflect.TypeOf(i).Elem()

	selectExpr := ""
	for i := 0; i < reflectedType.NumField(); i++ {
		field := reflectedType.Field(i)

		column := field.Tag.Get("column")
		if len(column) == 0 {
			continue
		}

		quotedColumn := fmt.Sprintf("`%s`", column)
		if len(selectExpr) > 0 {
			selectExpr += ", "
		}
		selectExpr += quotedColumn
	}

	return selectExpr
}

func (q *SelectQuery) Query(tx *sql.Tx, args ...interface{}) (interface{}, error) {
	if !q.tableSet {
		return nil, errors.WithStack(NoTablesUsed)
	}

	rows, err := database.Query(q.query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	result, err := q.iterate(rows)
	if err != nil {
		return nil, err
	}

	return result.Interface(), nil
}

func (q *SelectQuery) QueryRow(tx *sql.Tx, args ...interface{}) (interface{}, error) {
	if !q.tableSet {
		return nil, errors.WithStack(NoTablesUsed)
	}

	row := database.QueryRow(q.query, args...)

	result, err := q.scan(row)
	if err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return result.Addr().Interface(), nil
}

func (q *SelectQuery) iterate(r rows) (reflect.Value, error) {
	sliceType := reflect.SliceOf(q.t)
	resultValue := reflect.MakeSlice(sliceType, 0, 0)

	for r.Next() {
		rowValue, err := q.scan(r)
		if err != nil {
			return resultValue, errors.WithStack(err)
		}

		resultValue = reflect.Append(resultValue, rowValue)
	}

	if err := r.Err(); err != nil {
		return resultValue, errors.WithStack(err)
	}

	return resultValue, nil
}

func (q *SelectQuery) scan(r row) (reflect.Value, error) {
	rowValue := reflect.New(q.t)
	rowValue = rowValue.Elem()

	fields := []interface{}{}
	for i := 0; i < rowValue.NumField(); i++ {
		fieldValue := rowValue.Field(i)
		fieldType := q.t.Field(i)

		column := fieldType.Tag.Get("column")
		if len(column) == 0 {
			continue
		}

		fieldValue = fieldValue.Addr()
		ptr := fieldValue.Interface()
		fields = append(fields, ptr)
	}
	err := r.Scan(fields...)
	return rowValue, errors.WithStack(err)
}

func (q *SelectQuery) Error() error {
	return q.err
}
