package main

import (
	"github.com/01-edu/z01"
)

type point struct {
	x int
	y int
}

func setPoint(ptr *point) {
	ptr.x = 42
	ptr.y = 21
}

func printStr(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
}

func printInt(a int) {
	r := '0'
	if a/10 > 0 {
		printInt(a / 10)
	}
	for i := 0; i < a%10; i++ {
		r++
	}
	z01.PrintRune(r)
}

func main() {
	points := &point{}
	setPoint(points)
	printStr("x = ")
	printInt(points.x)
	printStr(", y = ")
	printInt(points.y)
	z01.PrintRune('\n')
}
