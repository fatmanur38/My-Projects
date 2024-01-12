package piscine

func ActiveBits(n int) int {
	count := 0
	for n != 0 {
		n = n & (n - 1)
		count++
	}
	return count
}
