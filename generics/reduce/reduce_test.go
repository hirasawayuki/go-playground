package main

import (
	"testing"
)

func TestReduce_IntSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "Returns the sum of the elements of a slice",
			input: []int{1, 2, 3, 4, 5},
			want:  15,
		},
	}

	for _, tt := range tests {
		got := reduce(tt.input, func(i, initValue int) int {
			return i + initValue
		}, 0)

		if got != tt.want {
			t.Errorf("got %v want %v\n", got, tt.want)
		}
	}
}

func TestReduce_StringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "Returns concatenate a string of the elements of a slice",
			input: []string{"G", "o", "l", "a", "n", "g"},
			want:  "Golang",
		},
	}

	for _, tt := range tests {
		got := reduce(tt.input, func(s, initValue string) string {
			return s + initValue
		}, "")

		if got != tt.want {
			t.Errorf("got %v want %v\n", got, tt.want)
		}
	}
}
