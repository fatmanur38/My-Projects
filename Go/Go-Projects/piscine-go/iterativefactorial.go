package piscine

func IterativeFactorial(nb int) int {
	faktoriyel := 1
	if nb < 0 || nb > 21 {
		return 0
	}
	for i := 1; i <= nb; i++ {
		faktoriyel *= i
	}
	return faktoriyel
}
