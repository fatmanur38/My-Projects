package piscine

import "github.com/01-edu/z01"

func PrintNbr(n int) {
	if n <= -9223372036854775808 {
		z01.PrintRune('0')
	}
	if n > 9223372036854775807 {
		z01.PrintRune('0')
	}
	if n < 0 {
		z01.PrintRune('-')
		n = -n
	}

	if n > 9 {
		PrintNbr(n / 10)
		PrintNbr(n % 10)
	} else {
		z01.PrintRune(rune(n + 48))
	}

	//if n == 0 {
	//z01.PrintRune('0')
	//}
}
