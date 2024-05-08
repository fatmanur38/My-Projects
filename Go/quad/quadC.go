package piscine

/*ABBBA     ABBBA
  B   B
  CBBBC*/
import "github.com/01-edu/z01"

func QuadC(x, y int) {
	for h := 1; h <= y; h++ {
		for w := 1; w <= x; w++ {
			if h == 1 {
				if w == 1 || w == x {
					z01.PrintRune('A')
				} else {
					z01.PrintRune('B')
				}
			} else if h < y {
				if w == 1 || w == x {
					z01.PrintRune('B')
				} else {
					z01.PrintRune(' ')
				}
			} else {
				if w == 1 || w == x {
					z01.PrintRune('C')
				} else {
					z01.PrintRune('B')
				}
			}
		}
		z01.PrintRune('\n')
	}
}
