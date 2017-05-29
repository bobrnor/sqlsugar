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
		query: "SELECT  FROM `EmptyTable`",
	}
	found := Select((*EmptyTable)(nil)).From([]string{"EmptyTable"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery1(t *testing.T) {
	expected := &SelectQuery{
		query: "SELECT `id`, `field0`, `field1` FROM `SimpleTable`",
	}
	found := Select((*SimpleTable)(nil)).From([]string{"SimpleTable"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

// func TestSelectQuery2(t *testing.T) {
// 	table1 := SimpleTable{
// 		TableMeta{
// 			TableName: "SimpleTable1",
// 		},
// 		127,
// 		"test field",
// 		3.14,
// 	}
// 	table2 := SimpleTable{
// 		TableMeta{
// 			TableName: "SimpleTable2",
// 		},
// 		512,
// 		"test field N2",
// 		6.626,
// 	}
// 	expected := &query{
// 		query: "SELECT `SimpleTable1`.`id`, `SimpleTable1`.`field0`, `SimpleTable1`.`field1`, `SimpleTable2`.`id`, `SimpleTable2`.`field0`, `SimpleTable2`.`field1` FROM `SimpleTable1`, `SimpleTable2`",
// 	}
// 	found, err := SelectQuery("SELECT @fields FROM `SimpleTable1`, `SimpleTable2`", &table1, &table2)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if !reflect.DeepEqual(expected, found) {
// 		t.Errorf("Expected: %+v, found %+v", expected, found)
// 	}
// }

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
