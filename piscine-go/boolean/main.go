package main

import (
	"os"

	"github.com/01-edu/z01"
)

const (
	EvenMsg = "I have an even number of arguments"
	OddMsg  = "I have an odd number of arguments"
)

type boolean bool

const (
	no  boolean = false
	yes boolean = true
)

func printStr(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
	z01.PrintRune('\n')
}

func even(nbr int) boolean {
	return boolean(nbr%2 == 0)
}

func isEven(nbr int) boolean {
	if even(nbr) == yes {
		return yes
	}
	return no
}

func main() {
	args := os.Args[1:]

	lengthOfArg := len(args)

	if isEven(lengthOfArg) == yes {
		printStr(EvenMsg)
	} else {
		printStr(OddMsg)
	}
}
