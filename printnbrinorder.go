package piscine

import (
	"github.com/01-edu/z01"
)

func PrintNbrInOrder(n int) {
	if n == 0 {
		z01.PrintRune('0')
	} else {
		count := 0
		sayi := n
		for n > 0 {
			n /= 10
			count++
		}
		rakamlar := make([]int, count)

		for i := count - 1; sayi > 0; i-- {
			rakamlar[i] = sayi % 10
			sayi /= 10
		}

		SortIntegerTable_1(rakamlar)

		for j := 0; j < len(rakamlar); j++ {
			PrintNbr_1(rakamlar[j])
		}
	}
}

func SortIntegerTable_1(table []int) {
	for i := 0; i < len(table); i++ {
		for j := i + 1; j < len(table); j++ {
			if table[i] > table[j] {
				gecici := table[j]
				table[j] = table[i]
				table[i] = gecici
			}
		}
	}
}

func PrintNbr_1(n int) {
	if n == 0 {
		z01.PrintRune('0')
		return
	}
	if n < 0 {
		z01.PrintRune('-')
		n = -n
	}
	num := uint(n)

	var rakamlar string
	for num != 0 {
		son_rakam := num % 10
		num /= 10
		rakamlar = string(rune('0'+son_rakam)) + rakamlar
	}

	for _, j := range rakamlar {
		z01.PrintRune(j)
	}
}
