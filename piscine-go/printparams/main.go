package main

import (
	"os"

	"github.com/01-edu/z01"
)

/*choumi
is
the
best
cat*/
func main() {
	fatma := os.Args
	for i := 1; i < len(os.Args); i++ {
		for _, a := range []rune(fatma[i]) {
			z01.PrintRune(a)
		}
		z01.PrintRune('\n')
	}

	/*var progekle string
	for i := 1; i < len(os.Args); i++ {
		progekle = os.Args[i]
		for _, char := range progekle {
			z01.PrintRune(char)
		}
		if i != len(os.Args) {
			z01.PrintRune('\n')
		}
	}
	z01.PrintRune('\n')*/

	/*var programekle string

	for i := 1; i < len(os.Args); i++ {
		programekle = os.Args[i]

		for _, char := range programekle {
			z01.PrintRune(char)
		}

		if i != len(os.Args)-1 {
			z01.PrintRune('\n')
		}
	}

	z01.PrintRune('\n')*/
}
