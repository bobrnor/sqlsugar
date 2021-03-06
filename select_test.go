package sqlsugar

import (
	"reflect"
	"testing"
)

type EmptyTable struct {
}

type SimpleTable struct {
	ID     int64   `column:"id"`
	Field0 string  `column:"field0"`
	Field1 float64 `column:"field1"`
}

func TestSelectQuery0(t *testing.T) {
	expected := &SelectQuery{
		query:    "SELECT  FROM `EmptyTable`",
		t:        reflect.TypeOf((*EmptyTable)(nil)).Elem(),
		tableSet: true,
	}
	found := Select((*EmptyTable)(nil)).From([]string{"EmptyTable"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery1(t *testing.T) {
	expected := &SelectQuery{
		query:    "SELECT `id`, `field0`, `field1` FROM `SimpleTable`",
		t:        reflect.TypeOf((*SimpleTable)(nil)).Elem(),
		tableSet: true,
	}
	found := Select((*SimpleTable)(nil)).From([]string{"SimpleTable"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery2(t *testing.T) {
	expected := &SelectQuery{
		query:    "SELECT `id`, `field0`, `field1` FROM `SimpleTable0`, `SimpleTable1`",
		t:        reflect.TypeOf((*SimpleTable)(nil)).Elem(),
		tableSet: true,
	}
	found := Select((*SimpleTable)(nil)).From([]string{"SimpleTable0", "SimpleTable1"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectExpression0(t *testing.T) {
	expected := ""
	found := selectExpression((*EmptyTable)(nil))

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectExpression1(t *testing.T) {
	expected := "`id`, `field0`, `field1`"
	found := selectExpression((*SimpleTable)(nil))

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

type Rows struct {
	values  []interface{}
	hasNext bool
}

func (r *Rows) Scan(dst ...interface{}) error {
	for i, v := range dst {
		value := reflect.ValueOf(r.values[i])
		reflect.ValueOf(v).Elem().Set(value)
	}
	return nil
}

func (r *Rows) Next() bool {
	if r.hasNext {
		r.hasNext = false
		return true
	}
	return false
}

func (r *Rows) Err() error {
	return nil
}

func TestScan0(t *testing.T) {
	r := Rows{
		values:  []interface{}{int64(127), "TEST_CLIENT_ID", float64(67.223)},
		hasNext: false,
	}
	expected := SimpleTable{
		ID:     int64(127),
		Field0: "TEST_CLIENT_ID",
		Field1: float64(67.223),
	}
	q := SelectQuery{
		t: reflect.TypeOf((*SimpleTable)(nil)).Elem(),
	}
	foundValue, err := q.scan(&r)
	found := foundValue.Interface().(SimpleTable)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%+v %+v", found, err)
	}
}

func TestIterate0(t *testing.T) {
	r := Rows{
		values:  []interface{}{int64(127), "TEST_CLIENT_ID", float64(67.223)},
		hasNext: true,
	}
	expected := []SimpleTable{
		{
			ID:     int64(127),
			Field0: "TEST_CLIENT_ID",
			Field1: float64(67.223),
		},
	}
	q := SelectQuery{
		t: reflect.TypeOf((*SimpleTable)(nil)).Elem(),
	}
	foundValue, err := q.iterate(&r)
	found := foundValue.Interface().([]SimpleTable)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%+v %+v", found, err)
	}
}
