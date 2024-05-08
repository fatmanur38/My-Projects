package piscine

// AtoiBase fonksiyonu, verilen bir string'i belirli bir sayı tabanına göre integer'a çevirir.
func AtoiBase(s string, base string) int {
	length := len(base)
	toplam := 0
	count := 0
	n := []int{}

	// Hatalı durumu kontrol et: Tabanın negatif olup olmadığını kontrol et
	if base[0] == '-' || s == "" {
		return 0
	}

	// Her karakteri tabanda karşılık gelen indekse dönüştür
	for i := 0; i < len(s); i++ {
		found := false
		for j := 0; j < length; j++ {
			if s[i] == base[j] {
				n = append(n, j)
				found = true
				break
			}
		}
		// Hatalı durumu kontrol et: Girişte geçersiz bir karakter varsa
		if !found {
			return 0
		}
	}

	// Belirli bir sayı tabanına göre integer'a çevir
	for i := 1; ; i *= length {
		if len(n)-1-count < 0 {
			break
		}
		toplam += n[len(n)-1-count] * i
		count++
	}

	return toplam
}
