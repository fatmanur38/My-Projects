package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arguments := os.Args[1:]
	for i := 0; i < len(arguments)-1; i++ {
		for j := 0; j < len(arguments)-i-1; j++ {
			if arguments[j] > arguments[j+1] {
				arguments[j], arguments[j+1] = arguments[j+1], arguments[j]
			}
		}
	}
	for i := 0; i < len(arguments); i++ {
		for _, value := range arguments[i] {
			z01.PrintRune(rune(value))
		}
		z01.PrintRune('\n')
	}
}
