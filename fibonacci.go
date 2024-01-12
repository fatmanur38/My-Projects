package piscine

func Fibonacci(index int) int {
	sonuc := 0
	if index < 0 {
		return -1
	} else if index == 0 {
		return 0
	} else if index == 1 {
		return 1
	} else {
		sonuc = Fibonacci(index-1) + Fibonacci(index-2)
	}
	return sonuc
}
