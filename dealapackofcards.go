package piscine

import (
	"fmt"

	"github.com/01-edu/z01"
)

func DealAPackOfCards(deck []int) {
	players := 4
	for i := 0; i < players; i++ {
		fmt.Printf("Player %d: ", i+1)
		fmt.Printf("%d", deck[i*3])
		z01.PrintRune(',')
		z01.PrintRune(' ')
		fmt.Printf("%d", deck[i*3+1])
		z01.PrintRune(',')
		z01.PrintRune(' ')
		fmt.Printf("%d", deck[i*3+2])
		z01.PrintRune('\n')
	}
}
