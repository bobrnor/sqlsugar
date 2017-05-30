package sqlsugar

import (
	"reflect"
	"testing"
)

func TestDeleteQuery0(t *testing.T) {
	expected := &DeleteQuery{
		query: "DELETE FROM `EmptyTable`",
	}
	found := Delete("EmptyTable")

	if !reflect.DeepEqual(expected, found) {
		t.Errorf("Expected: %+v, found %+v", expected, found)
	}
	t.Errorf("Expected: %+v, found %+v", expected, found)
}
