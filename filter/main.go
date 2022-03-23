package main

import "fmt"

func filter[T any](s []T, f func(T) bool) []T {
	n := make([]T, 0, len(s))
	for _, e := range s {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func main() {
	chars := []string{"a", "b", "c", "d"}
	c := filter(chars, func(s string) bool {
		return s == "a" || s == "b"
	})
	fmt.Println(c)

	numbers := []int{1, 2, 3, 4, 5, 6}
	n := filter(numbers, func(i int) bool {
		return i%2 == 0
	})
	fmt.Println(n)
}
