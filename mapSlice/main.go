package main

import (
	"fmt"
)

func mapSlice[T, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
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
