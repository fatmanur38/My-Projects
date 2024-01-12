package piscine

func Atoi(s string) int {
	result := 0
	sign := 1

	// Eğer stringin başında + veya - işareti varsa, işareti kontrol et ve ilerlemeye devam et.
	if s != "" && (s[0] == '+' || s[0] == '-') {
		if s[0] == '-' {
			sign = -1
		}
		s = s[1:] // İlk karakteri atla
	}

	for _, char := range s {
		if char >= '0' && char <= '9' {
			digit := int(char - '0')
			result = result*10 + digit
		} else {
			// Geçersiz karakter bulunduğunda hemen 0 döndür.
			return 0
		}
	}

	return result * sign
}
