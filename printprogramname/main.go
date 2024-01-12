package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	fatma := os.Args[0]
	necla := []rune(fatma)
	for _, a := range necla[2:] {
		z01.PrintRune(a)
	}
	z01.PrintRune(('\n'))
}
