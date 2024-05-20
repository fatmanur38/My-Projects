package piscine

func IsLower(s string) bool {
	array := []rune(s)
	for i := 0; i < len(s); i++ {
		if !(array[i] >= 'a' && array[i] <= 'z') {
			return false
		}
	}
	return true
}
