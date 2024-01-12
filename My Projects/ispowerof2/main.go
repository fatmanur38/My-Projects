package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arg := os.Args[1:]
	if len(arg) != 1 {
		return
	}
	nbr := Atoi(arg[0])

	if ispowerof(nbr) {
		printstr("true\n")
	} else {
		printstr("false\n")
	}
}

func Atoi(s string) int {
	isaret := 1
	result := 0
	for index, i := range s {
		if index == 0 && i == '+' {
			isaret = 1
		} else if index == 0 && i == '-' {
			isaret = -1
		} else if i >= '0' && i <= '9' {
			result = result*10 + int(i-48)
		} else {
			return 0
		}
	}
	return result * isaret
}

func ispowerof(n int) bool {
	if n == 1 {
		return true
	} else if n%2 != 0 {
		return false
	}
	return ispowerof(n / 2)
}

func printstr(s string) {
	for _, i := range s {
		z01.PrintRune(i)
	}
}
