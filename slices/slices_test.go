package slices

import (
	"reflect"
	"testing"
)

func TestReverseInPlace(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
	}

	for _, test := range tests {
		Reverse(test.input)
		if !(reflect.DeepEqual(test.input, test.want)) {
			t.Errorf("wanted: %v, got: %v", test.want, test.input)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
	}

	for _, test := range tests {
		got := Reversed(test.input)
		if !(reflect.DeepEqual(got, test.want)) {
			t.Errorf("wanted: %v, got: %v", test.want, got)
		}
	}
}
