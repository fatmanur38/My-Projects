package piscine

import (
	"github.com/01-edu/z01"
)

// CheckBaseValidity, verilen base'in geçerli olup olmadığını kontrol eder.
func CheckBaseValidity(base string) bool {
	// Bir base en az 2 karakter içermelidir.
	if len(base) < 2 {
		return false
	}

	// Her bir karakterin benzersiz olup olmadığını kontrol eder.
	seen := make(map[rune]bool)
	for _, char := range base {
		if char == '+' || char == '-' {
			return false
		}
		if seen[char] {
			return false
		}
		seen[char] = true
	}

	return true
}

// PrintNbrBase, verilen sayıyı belirtilen tabanda yazdıran fonksiyon.
func PrintNbrBase(nbr int, base string) {
	if !CheckBaseValidity(base) {
		// Geçersiz bir base ise "NV" yazdır.
		z01.PrintRune('N')
		z01.PrintRune('V')
		return
	}

	// Sayı sıfır veya negatif ise "-" yazdır ve pozitife çevir.
	if nbr < 0 {
		z01.PrintRune('-')
		nbr *= -1
	}

	// Sayıyı belirtilen tabanda yazdır.
	baseLen := len(base)
	if nbr >= baseLen {
		PrintNbrBase(nbr/baseLen, base)
	}
	z01.PrintRune(rune(base[nbr%baseLen]))
}
