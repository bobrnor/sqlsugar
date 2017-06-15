package sqlsugar

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestUpdateQuery0(t *testing.T) {
	expected := &UpdateQuery{
		err: NoSet,
	}
	found := Update("EmptyTable").Set([]string{})

	if !reflect.DeepEqual(expected.err, errors.Cause(found.err)) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery1(t *testing.T) {
	expected := &UpdateQuery{
		query:      "UPDATE `SimpleTable` SET `id` = ?, `field0` = ?, `field1` = ?",
		setColumns: []string{"id", "field0", "field1"},
	}
	found := Update("SimpleTable").Set([]string{"id", "field0", "field1"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery2(t *testing.T) {
	expected := &UpdateQuery{
		query:      "UPDATE `SimpleTable` SET `field0` = ?",
		setColumns: []string{"field0"},
	}
	found := Update("SimpleTable").Set([]string{"field0"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestUpdateQuery3(t *testing.T) {
	var table struct {
		TestField string `column:"test_field"`
	}
	table.TestField = "test"
	found := UpdateMultiple([]string{"SimpleTable1", "SimpleTable2"}).SetAll(&table)

	if errors.Cause(found.Error()) != InappropriateSetAllUsage {
		t.Errorf("Expected error: %+v, found %+v", InappropriateSetAllUsage, found)
	}
}

func TestUpdateQuery4(t *testing.T) {
	expected := &UpdateQuery{
		query:          "UPDATE `SimpleTable1`, `SimpleTable2` SET `SimpleTable2`.`field0` = ?",
		setColumns:     []string{"field0"},
		multipleTables: true,
	}
	found := UpdateMultiple([]string{"SimpleTable1", "SimpleTable2"}).Set([]string{"SimpleTable2.field0"})

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}
