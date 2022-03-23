package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	got := filter(numbers, func(i int) bool {
		return i%2 == 0
	})
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("filter got %v want %v\n", got, want)
	}

	chars := []string{"apple", "banana", "cat", "dog", "gorilla"}
	got2 := filter(chars, func(str string) bool {
		return strings.Contains(str, "o")
	})
	want2 := []string{"dog", "gorilla"}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("filter got %v want %v\n", got, want)
	}
}
