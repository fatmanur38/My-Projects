package piscine

import "github.com/01-edu/z01"

/*
ABBBC

	B   B
	CBBBA
*/
func QuadE(x, y int) {
	for h := 1; h <= y; h++ {
		for w := 1; w <= x; w++ {
			if h == 1 {
				if w == 1 {
					z01.PrintRune('A')
				} else if w < x {
					z01.PrintRune('B')
				} else {
					z01.PrintRune('C')
				}
			} else if h < y {
				if w == 1 || w == x {
					z01.PrintRune('B')
				} else {
					z01.PrintRune(' ')
				}
			} else {
				if w == 1 {
					z01.PrintRune('C')
				} else if w < x {
					z01.PrintRune('B')
				} else {
					z01.PrintRune('A')
				}
			}
		}
		z01.PrintRune('\n')
	}
}
