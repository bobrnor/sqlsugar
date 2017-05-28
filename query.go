package sorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type TableMeta struct {
	TableName string
}

type Query struct {
	query string
}

type row interface {
	Scan(dest ...interface{}) error
}

type rows interface {
	row
	Next() bool
	Err() error
}

const (
	SelectExprPlaceholder = "@fields"

	UpdateTableExprPlaceholder = "@table"
	UpdateSetsExprPlaceholder  = "@sets"

	DeleteTableExprPlaceholder = "@table"
)

var (
	MetaNotFound = errors.New("Table struct does not contain meta info")
)

func SelectQuery(format string, args ...interface{}) (*Query, error) {
	jointSelectExpr := ""

	for _, i := range args {
		selectExpr, err := selectExpression(i)
		if err != nil {
			return nil, err
		}

		if len(jointSelectExpr) > 0 {
			jointSelectExpr += ", "
		}
		jointSelectExpr += selectExpr
	}

	return &Query{
		query: strings.Replace(format, SelectExprPlaceholder, jointSelectExpr, -1),
	}, nil
}

func selectExpression(i interface{}) (string, error) {
	meta, err := fetchMeta(i)
	if err != nil {
		return "", err
	}

	reflectedType := reflect.TypeOf(i).Elem()

	selectExpr := ""
	for i := 0; i < reflectedType.NumField(); i++ {
		field := reflectedType.Field(i)

		column := field.Tag.Get("column")
		if len(column) == 0 {
			continue
		}

		quotedColumn := fmt.Sprintf("`%s`.`%s`", meta.TableName, column)
		if len(selectExpr) > 0 {
			selectExpr += ", "
		}
		selectExpr += quotedColumn
	}

	return selectExpr, nil
}

func InsertQuery(i interface{}) (*Query, error) {
	meta, err := fetchMeta(i)
	if err != nil {
		return nil, err
	}

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
		query: fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", meta.TableName, columnsExpr, valuesExpr),
	}, nil
}

func UpdateQuery(format string, fields []string, i interface{}) (*Query, error) {
	meta, err := fetchMeta(i)
	if err != nil {
		return nil, err
	}

	reflectedType := reflect.TypeOf(i).Elem()

	setsExpr := ""
	for i := 0; i < reflectedType.NumField(); i++ {
		field := reflectedType.Field(i)

		column := field.Tag.Get("column")
		if !contains(fields, column) {
			continue
		}

		if len(setsExpr) > 0 {
			setsExpr += ", "
		}

		setsExpr += fmt.Sprintf("`%s` = ?", column)
	}

	quotedTableName := fmt.Sprintf("`%s`", meta.TableName)
	update := strings.Replace(format, UpdateTableExprPlaceholder, quotedTableName, -1)
	update = strings.Replace(update, UpdateSetsExprPlaceholder, setsExpr, -1)
	return &Query{
		query: update,
	}, nil
}

func DeleteQuery(format string, i interface{}) (*Query, error) {
	meta, err := fetchMeta(i)
	if err != nil {
		return nil, err
	}

	quotedTableName := fmt.Sprintf("`%s`", meta.TableName)
	return &Query{
		query: strings.Replace(format, DeleteTableExprPlaceholder, quotedTableName, -1),
	}, nil
}

func fetchMeta(i interface{}) (TableMeta, error) {
	value := reflect.ValueOf(i).Elem()
	metaValue := value.FieldByName("TableMeta").Interface()
	meta, ok := metaValue.(TableMeta)
	if !ok {
		return TableMeta{}, errors.WithStack(MetaNotFound)
	}

	return meta, nil
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
