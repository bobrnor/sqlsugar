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
	expected := &Query{
		query: "SELECT  FROM `EmptyTable`",
	}
	found := Select((*EmptyTable)(nil)).From([]string{"EmptyTable"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery1(t *testing.T) {
	expected := &Query{
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
// 	expected := &Query{
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
	found, err := selectExpression((*EmptyTable)(nil))

	if err != nil {
		t.Error(err)
	}

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectExpression1(t *testing.T) {
	expected := "`id`, `field0`, `field1`"
	found, err := selectExpression((*SimpleTable)(nil))

	if err != nil {
		t.Error(err)
	}

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestInsertQuery0(t *testing.T) {
	expected := &Query{
		query: "INSERT INTO `EmptyTable` () VALUES ()",
	}
	found := Insert((*EmptyTable)(nil)).Into("EmptyTable")

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestInsertQuery1(t *testing.T) {
	expected := &Query{
		query: "INSERT INTO `SimpleTable` (`field0`, `field1`) VALUES (?, ?)",
	}
	found := Insert((*SimpleTable)(nil)).Into("SimpleTable")

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery0(t *testing.T) {
	expected := &Query{
		err: NoSetColumns,
	}
	found := Update("EmptyTable").Set([]string{})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery1(t *testing.T) {
	expected := &Query{
		query: "UPDATE `SimpleTable` SET `id` = ?, `field0` = ?, `field1` = ?",
	}
	found := Update("SimpleTable").Set([]string{"id", "field0", "field1"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery2(t *testing.T) {
	expected := &Query{
		query: "UPDATE `SimpleTable` SET `field0` = ?",
	}
	found := Update("SimpleTable").Set([]string{"field0"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestDeleteQuery0(t *testing.T) {
	expected := &Query{
		query: "DELETE FROM `EmptyTable`",
	}
	found := Delete("EmptyTable")

	if !reflect.DeepEqual(expected, found) {
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
	found := SimpleTable{}
	foundType := reflect.TypeOf(found)
	foundValue, err := scan(&r, foundType)
	found = foundValue.Interface().(SimpleTable)
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
	found := []SimpleTable{}
	foundType := reflect.TypeOf(found)
	foundValue, err := iterate(&r, foundType)
	found = foundValue.Interface().([]SimpleTable)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%+v %+v", found, err)
	}
}
