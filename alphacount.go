package piscine

//"Hello 78 World!    4455 /" == 10
func AlphaCount(s string) int {
	array := []rune(s)
	harf := 0
	for i := 0; i < len(s); i++ {
		if (array[i] >= 'a' && array[i] <= 'z') || array[i] >= 'A' && array[i] <= 'Z' {
			harf++
		}
	}
	return harf
}
