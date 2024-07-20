package main

import "fmt"

func main() {
	for c := range Alphabet {
		fmt.Printf("%c", c)
		if c == 'C' {
			break
		}
	}
}

func Alphabet(yield func(rune) bool) {
	for c := 'A'; c <= 'Z'; c++ {
		if !yield(c) {
			break
		}
	}
}
