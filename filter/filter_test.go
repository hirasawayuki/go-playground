package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilter_IntSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "filter even numbers",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{2, 4},
		}, {
			name:  "filter even numbers (empty slice)",
			input: []int{},
			want:  []int{},
		}, {
			name:  "filter even numbers (no match)",
			input: []int{1, 3, 5, 7, 9},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		got := filter(tt.input, func(i int) bool {
			return i%2 == 0
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("filter got %v want %v\n", got, tt.want)
		}
	}
}

func TestFilter_StringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "filter contains keyword",
			input: []string{"apple", "banana", "cat", "dog", "gorilla"},
			want:  []string{"dog", "gorilla"},
		}, {
			name:  "filter contains keyword (empty slice)",
			input: []string{},
			want:  []string{},
		}, {
			name:  "filter contains keyword (no match keyword)",
			input: []string{"a", "b", "c", "d", "e"},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		got := filter(tt.input, func(str string) bool {
			return strings.Contains(str, "o")
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("filter got %v want %v\n", got, tt.want)
		}
	}
}
