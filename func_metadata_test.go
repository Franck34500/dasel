package dasel

import (
	"reflect"
	"testing"
)

func TestMetadataFunc(t *testing.T) {
	t.Run("Type", func(t *testing.T) {
		orig := []interface{}{
			"abc", true, false, 1, 1.1, []interface{}{1},
		}
		ctx := NewContext(&orig, "all().metadata(type)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		exp := []interface{}{
			"string", "bool", "bool", "int", "float64", "slice",
		}
		got := s.Interfaces()

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})
}