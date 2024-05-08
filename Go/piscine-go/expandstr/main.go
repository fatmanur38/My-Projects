package main

import (
	"os"
	"strings"

	"github.com/01-edu/z01"
)

func main() {
	if len(os.Args) != 2 {
		return
	}
	words := strings.Fields(os.Args[1])
	var expand string
	for i, word := range words {
		expand += word
		if i != len(words)-1 {
			expand += "   "
		}
	}
	for _, r := range expand {
		z01.PrintRune(r)
	}
	z01.PrintRune('\n')
}
