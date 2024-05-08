package piscine

func IsUpper(s string) bool {
	array := []rune(s)
	for i := 0; i < len(s); i++ {
		if !(array[i] >= 'A' && array[i] <= 'Z') {
			return false
		}
	}
	return true
}
