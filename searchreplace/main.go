package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	if len(os.Args) == 4 {
		x := os.Args[1]
		y := os.Args[2]
		z := os.Args[3]
		result := ""
		for _, r := range x {
			if string(r) == y {
				result += z
			} else {
				result += string(r)
			}
		}
		for _, r := range result {
			z01.PrintRune(r)
		}
		z01.PrintRune('\n')
	} else {
		return
	}
}
