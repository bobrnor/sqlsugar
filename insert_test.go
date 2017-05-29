package sqlsugar

import (
	"reflect"
	"testing"
)

func TestInsertQuery0(t *testing.T) {
	expected := &InsertQuery{
		query: "INSERT INTO `EmptyTable` () VALUES ()",
	}
	found := Insert((*EmptyTable)(nil)).Into("EmptyTable")

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}

func TestInsertQuery1(t *testing.T) {
	expected := &InsertQuery{
		query: "INSERT INTO `SimpleTable` (`field0`, `field1`) VALUES (?, ?)",
	}
	found := Insert((*SimpleTable)(nil)).Into("SimpleTable")

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
}
