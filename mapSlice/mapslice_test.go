package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestMapSlice_IntSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "Returns a slice of the element multiplied by 2",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{2, 4, 6, 8, 10},
		},
	}

	for _, tt := range tests {
		got := mapSlice(tt.input, func(i int) int {
			return i * 2
		})

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v\n", got, tt.want)
		}
	}
}

func TestMapSlice_StringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "Returns a slice with the element capitalized",
			input: []string{"a", "b", "c", "d", "e"},
			want:  []string{"A", "B", "C", "D", "E"},
		},
	}

	for _, tt := range tests {
		got := mapSlice(tt.input, func(s string) string {
			return strings.ToUpper(s)
		})

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v\n", got, tt.want)
		}
	}
}
