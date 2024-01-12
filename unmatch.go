package piscine

// Unmatch fonksiyonu, bir dilimde çift olmayan elemanı döndürür.
// Eğer tüm sayılar bir çiftle eşleşiyorsa, -1 döndürür.
func Unmatch(a []int) int {
	// Dilimi tarayarak her bir elemanın sayısını güncelle
	for i := 0; i < len(a); i++ {
		count := 0
		for j := 0; j < len(a); j++ {
			// Eğer a[i] ile eşleşen bir eleman bulunduysa, sayacı artır
			if a[i] == a[j] {
				count++
			}
		}
		// Sayacın değeri tekse, yani çift sayıda eşleşme yoksa, bu elemanı döndür
		if count%2 != 0 {
			return a[i]
		}
	}

	// Eğer tüm sayılar bir çiftle eşleşiyorsa, -1 döndür
	return -1
}
