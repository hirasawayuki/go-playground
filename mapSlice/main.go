package main

import (
	"fmt"
)

func mapSlice[T, U any](sl []T, f func(T) U) []U {
	r := make([]U, len(sl))
	for i, s := range sl {
		r[i] = f(s)
	}
	return r
}

func main() {
	numbers := []float64{4, 9, 16, 25}
	newNumbers := mapSlice(numbers, func(i float64) float64 {
		return i * 2
	})
	fmt.Println(newNumbers)

	words := []string{"a", "b", "c", "d"}
	quoted := mapSlice(words, func(s string) string {
		return "\"" + s + "\""
	})
	fmt.Println(quoted)
}
