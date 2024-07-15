package main

import "fmt"

func main() {
	var ss []string = Map([]int{10, 20}, func(n int) string {
		return fmt.Sprintf("0x%x", n)
	})

	fmt.Println(ss)
}

func Map[S, T any](s []S, f func(n S) T) []T {
	ss := make([]T, 0, len(s))

	for _, v := range s {
		ss = append(ss, f(v))
	}

	return ss
}
