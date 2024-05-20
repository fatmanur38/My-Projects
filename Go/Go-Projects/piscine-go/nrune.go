package piscine

func NRune(s string, n int) rune {
	length := len(s)
	istenen := []rune(s)
	if n <= 0 || n > length {
		return 0
	}
	return istenen[n-1]
}
