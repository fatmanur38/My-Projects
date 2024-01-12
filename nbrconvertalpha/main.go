package main

import (
	"os"
	"strconv"

	"github.com/01-edu/z01"
)

func main() {
	args := os.Args[1:]

	// --upper bayrağı mevcut mu diye kontrol et
	upper := false
	if len(args) > 0 && args[0] == "--upper" {
		upper = true
		args = args[1:]
	}

	// Her bir argümanı işle
	for _, arg := range args {
		// Argümanı bir tamsayıya çevir
		position, err := strconv.Atoi(arg)
		if err != nil {
			// Geçersiz argüman, boşluk karakteri yazdır
			z01.PrintRune(' ')
			continue
		}

		// Konumun geçerli aralıkta olup olmadığını kontrol et (1 ile 26 arasında)
		if position >= 1 && position <= 26 {
			// Karşılık gelen harfi hesapla
			letter := 'a' + rune(position-1)
			if upper {
				// --upper bayrağı mevcutsa büyük harfe çevir
				letter = 'A' + rune(position-1)
			}
			z01.PrintRune(letter)
		} else {
			// Geçersiz konum, boşluk karakteri yazdır
			z01.PrintRune(' ')
		}
	}

	// Sonunda yeni satırı yazdır
	z01.PrintRune('\n')
}
