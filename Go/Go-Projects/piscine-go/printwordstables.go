package piscine

import "github.com/01-edu/z01"

func PrintWordsTables(a []string) {
	var words []string
	hiclik := ""

	for _, karakter := range a {
		if karakter == " " || karakter == "\t" || karakter == "\n" {
			if hiclik != "" {
				words = append(words, hiclik)
				hiclik = ""
			}
		} else {
			hiclik += karakter
		}
	}

	if hiclik != "" {
		words = append(words, hiclik)
	}

	// Ayırdığımız kelimeleri yazdıralım.
	for i, word := range words {
		for _, char := range word {
			z01.PrintRune(char)
		}
		if i != len(words)-1 {
			z01.PrintRune('\n')
		}
	}
}
