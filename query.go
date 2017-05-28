package sqlsugar

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type Query struct {
	query string
	err   error
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
	NoFromTables = errors.New("Tables not found in `FROM` statement")
	NoIntoTable  = errors.New("No table in `INSERT INTO` statement")
	NoSetColumns = errors.New("Columns not  found in `SET` statement")
)

func Select(i interface{}) *Query {
	selectExpr, err := selectExpression(i)
	if err != nil {
		return &Query{
			err: err,
		}
	}

	return &Query{
		query: fmt.Sprintf("SELECT %s", selectExpr),
	}
}

func (q *Query) From(tables []string) *Query {
	if q.err != nil {
		return q
	}

	if len(tables) == 0 {
		return &Query{
			query: "",
			err:   NoFromTables,
		}
	}

	quotedTables := []string{}
	for _, table := range tables {
		quotedTables = append(quotedTables, fmt.Sprintf("`%s`", table))
	}

	return &Query{
		query: fmt.Sprintf("%s FROM %s", q.query, strings.Join(quotedTables, ", ")),
	}
}

func (q *Query) Where(condition string) *Query {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &Query{
		query: fmt.Sprintf("%s WHERE %s", q.query, condition),
	}
}

func (q *Query) GroupBy(condition string) *Query {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &Query{
		query: fmt.Sprintf("%s GROUP BY %s", q.query, condition),
	}
}

func (q *Query) Having(condition string) *Query {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &Query{
		query: fmt.Sprintf("%s HAVING %s", q.query, condition),
	}
}

func (q *Query) OrderBy(condition string) *Query {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &Query{
		query: fmt.Sprintf("%s ORDER BY %s", q.query, condition),
	}
}

func (q *Query) Limit(condition string) *Query {
	if q.err != nil {
		return q
	}

	if len(condition) == 0 {
		return q
	}

	return &Query{
		query: fmt.Sprintf("%s LIMIT %s", q.query, condition),
	}
}

func selectExpression(i interface{}) (string, error) {
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

	return selectExpr, nil
}

func Insert(i interface{}) *Query {
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
	return &Query{
		query: fmt.Sprintf("(%s) VALUES (%s)", columnsExpr, valuesExpr),
	}
}

func (q *Query) Into(table string) *Query {
	if q.err != nil {
		return q
	}

	if len(table) == 0 {
		return &Query{
			query: "",
			err:   NoIntoTable,
		}
	}

	return &Query{
		query: fmt.Sprintf("INSERT INTO `%s` %s", table, q.query),
	}
}

func Update(table string) *Query {
	return &Query{
		query: fmt.Sprintf("UPDATE `%s`", table),
	}
}

func (q *Query) Set(columns []string) *Query {
	if q.err != nil {
		return q
	}

	if len(columns) == 0 {
		return &Query{
			query: "",
			err:   NoSetColumns,
		}
	}

	setsExpr := ""
	for _, column := range columns {
		if len(setsExpr) > 0 {
			setsExpr += ", "
		}

		setsExpr += fmt.Sprintf("`%s` = ?", column)
	}

	return &Query{
		query: fmt.Sprintf("%s SET %s", q.query, setsExpr),
	}
}

func Delete(table string) *Query {
	return &Query{
		query: fmt.Sprintf("DELETE FROM `%s`", table),
	}
}

func contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func (q *Query) Query(tx *sql.Tx, args []interface{}, i interface{}) error {
	rows, err := database.Query(q.query, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	resultType := reflect.TypeOf(i).Elem()
	resultValue, err := iterate(rows, resultType)
	if err != nil {
		return err
	}

	iValue := reflect.ValueOf(i)
	iValue = iValue.Elem()
	iValue.Set(resultValue)

	return nil
}

func (q *Query) QueryRow(tx *sql.Tx, args []interface{}, i interface{}) error {
	row := database.QueryRow(q.query, args...)

	resultType := reflect.TypeOf(i).Elem()
	resultValue, err := scan(row, resultType)
	if err != nil {
		return err
	}

	iValue := reflect.ValueOf(i)
	iValue = iValue.Elem()
	iValue.Set(resultValue)

	return nil
}

func (q *Query) Exec(tx *sql.Tx, args []interface{}) (sql.Result, error) {
	result, err := database.Exec(q.query, args...)
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func iterate(r rows, resultType reflect.Type) (reflect.Value, error) {
	resultValue := reflect.Zero(resultType)

	rowType := resultType.Elem()
	for r.Next() {
		rowValue, err := scan(r, rowType)
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

func scan(r row, rowType reflect.Type) (reflect.Value, error) {
	rowValue := reflect.New(rowType)
	rowValue = rowValue.Elem()

	fields := []interface{}{}
	for i := 0; i < rowValue.NumField(); i++ {
		fieldValue := rowValue.Field(i)
		fieldType := rowType.Field(i)

		column := fieldType.Tag.Get("column")
		if len(column) == 0 {
			continue
		}

		fieldValue = fieldValue.Addr()
		ptr := fieldValue.Interface()
		fields = append(fields, ptr)
	}
	err := r.Scan(fields...)
	return rowValue, err
}

func (q *Query) Error() error {
	return q.err
}
