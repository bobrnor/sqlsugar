package sqlsugar

import (
	"reflect"
	"testing"
)

type EmptyTable struct {
	TableMeta
}

type SimpleTable struct {
	TableMeta
	ID     int64   `column:"id"`
	Field0 string  `column:"field0"`
	Field1 float64 `column:"field1"`
}

func TestSelectQuery0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	expected := &Query{
		query: "SELECT  FROM `EmptyTable`",
	}
	found, err := SelectQuery("SELECT @fields FROM `EmptyTable`", &table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery1(t *testing.T) {
	table := SimpleTable{
		TableMeta{
			TableName: "SimpleTable",
		},
		127,
		"test field",
		3.14,
	}
	expected := &Query{
		query: "SELECT `SimpleTable`.`id`, `SimpleTable`.`field0`, `SimpleTable`.`field1` FROM `SimpleTable`",
	}
	found, err := SelectQuery("SELECT @fields FROM `SimpleTable`", &table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectQuery2(t *testing.T) {
	table1 := SimpleTable{
		TableMeta{
			TableName: "SimpleTable1",
		},
		127,
		"test field",
		3.14,
	}
	table2 := SimpleTable{
		TableMeta{
			TableName: "SimpleTable2",
		},
		512,
		"test field N2",
		6.626,
	}
	expected := &Query{
		query: "SELECT `SimpleTable1`.`id`, `SimpleTable1`.`field0`, `SimpleTable1`.`field1`, `SimpleTable2`.`id`, `SimpleTable2`.`field0`, `SimpleTable2`.`field1` FROM `SimpleTable1`, `SimpleTable2`",
	}
	found, err := SelectQuery("SELECT @fields FROM `SimpleTable1`, `SimpleTable2`", &table1, &table2)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectExpression0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	expected := ""
	found, err := selectExpression(&table)

	if err != nil {
		t.Error(err)
	}

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestSelectExpression1(t *testing.T) {
	table := SimpleTable{
		TableMeta{
			TableName: "SimpleTable",
		},
		127,
		"test field",
		3.14,
	}
	expected := "`SimpleTable`.`id`, `SimpleTable`.`field0`, `SimpleTable`.`field1`"
	found, err := selectExpression(&table)

	if err != nil {
		t.Error(err)
	}

	if expected != found {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestFetchMeta0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	found, err := fetchMeta(&table)

	if err != nil {
		t.Error(err)
	}

	if table.TableMeta != found {
		t.Errorf("Expected: %+v, found %+v", table.TableMeta, found)
	}
}

func TestInsertQuery0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	expected := &Query{
		query: "INSERT INTO `EmptyTable` () VALUES ()",
	}
	found, err := InsertQuery(&table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestInsertQuery1(t *testing.T) {
	table := SimpleTable{
		TableMeta{
			TableName: "SimpleTable",
		},
		127,
		"test field",
		3.14,
	}
	expected := &Query{
		query: "INSERT INTO `SimpleTable` (`field0`, `field1`) VALUES (?, ?)",
	}
	found, err := InsertQuery(&table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	expected := &Query{
		query: "UPDATE `EmptyTable` SET ",
	}
	found, err := UpdateQuery("UPDATE @table SET @sets", []string{}, &table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery1(t *testing.T) {
	table := SimpleTable{
		TableMeta{
			TableName: "SimpleTable",
		},
		127,
		"test field",
		3.14,
	}
	expected := &Query{
		query: "UPDATE `SimpleTable` SET `id` = ?, `field0` = ?, `field1` = ?",
	}
	found, err := UpdateQuery("UPDATE @table SET @sets", []string{"id", "field0", "field1"}, &table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery2(t *testing.T) {
	table := SimpleTable{
		TableMeta{
			TableName: "SimpleTable",
		},
		127,
		"test field",
		3.14,
	}
	expected := &Query{
		query: "UPDATE `SimpleTable` SET `field0` = ?",
	}
	found, err := UpdateQuery("UPDATE @table SET @sets", []string{"field0"}, &table)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestDeleteQuery0(t *testing.T) {
	table := EmptyTable{
		TableMeta{
			TableName: "EmptyTable",
		},
	}
	expected := &Query{
		query: "DELETE FROM `EmptyTable`",
	}
	found, err := DeleteQuery("DELETE FROM @table", &table)

	if err != nil {
		t.Error(err)
	}

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
