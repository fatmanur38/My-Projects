package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		n := 0
		for _, val := range os.Args[1] {
			if val < '0' || val > '9' {
				fmt.Printf("ERROR: cannot convert to roman digit\n")
				return
			}
			n = n*10 + int(val-48)
		}
		if n < 1 || n > 3999 {
			fmt.Printf("ERROR: cannot convert to roman digit\n")
			return
		}
		num := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		sum := []string{"M", "(M-C)", "D", "(D-C)", "C", "(C-X)", "L", "(L-X)", "X", "(X-I)", "V", "(V-I)", "I"}
		roman := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
		var res1, res2 string
		for val := range num {
			tabl := n / num[val]
			for tabl > 0 {
				n -= num[val]
				tabl--
				res1 += sum[val] + "+"
				res2 += roman[val]
			}
		}
		fmt.Printf("%s\n%s\n", res1[:len(res1)-1], res2)
	}
}
