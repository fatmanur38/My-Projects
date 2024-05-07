package piscine

func BasicAtoi2(s string) int {
	result := 0

	for _, char := range s {
		// Eğer karakter bir rakam değilse, fonksiyon 0 döndürür.
		if char < '0' || char > '9' {
			return 0
		}

		// Her bir karakterin ASCII değeri '0' karakterinin ASCII değerinden çıkartılarak,
		// gerçek rakam değeri elde edilir ve bu değer sonuç değişkenine eklenir.
		digit := int(char - '0')
		result = result*10 + digit
	}

	return result
}
