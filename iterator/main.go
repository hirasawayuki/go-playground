package main

import (
	"fmt"
	"slices"
)

func main() {
	for c := range Alphabet {
		fmt.Printf("%c", c)
		if c == 'C' {
			break
		}
	}

	fmt.Println("")
	values := []string{"a", "b", "c"}

	for i, s := range slices.All(values) {
		fmt.Printf("%d: %s\n", i, s)
	}
}

func Alphabet(yield func(rune) bool) {
	for c := 'A'; c <= 'Z'; c++ {
		if !yield(c) {
			break
		}
	}
}
