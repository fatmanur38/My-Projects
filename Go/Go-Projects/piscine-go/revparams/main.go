package main

import (
	"os"

	"github.com/01-edu/z01"
)

/*cat
best
the
is
choumi*/
func main() {
	fatma := os.Args
	for i := len(fatma) - 1; i > 0; i-- {
		for _, a := range []rune(fatma[i]) {
			z01.PrintRune(a)
		}
		z01.PrintRune('\n')
	}
}
