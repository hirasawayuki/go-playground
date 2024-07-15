package main

import "fmt"

func main() {
	var sum int
	ss1 := []int{1, 2, 3, 4, 5}

	ss1 = Filter(ss1, func(n int) bool {
		return n%2 == 1
	})

	Apply(ss1, func(n int) {
		sum += n
	})

	fmt.Println(sum)

	var ss []string = Map([]int{10, 20}, func(n int) string {
		return fmt.Sprintf("0x%x", n)
	})

	fmt.Println(ss)
}

func Apply(s []int, f func(n int)) {
	for _, v := range s {
		f(v)
	}
}

func Filter[T any](s []T, f func(n T) bool) []T {
	var ss []T

	for _, v := range s {
		if f(v) {
			ss = append(ss, v)
		}
	}

	return ss
}

func Map[S, T any](s []S, f func(n S) T) []T {
	ss := make([]T, 0, len(s))

	for _, v := range s {
		ss = append(ss, f(v))
	}

	return ss
}
