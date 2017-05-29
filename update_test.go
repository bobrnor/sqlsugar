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
