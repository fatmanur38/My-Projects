package piscine

func IsNumeric(s string) bool {
	array := []rune(s)
	for i := 0; i < len(s); i++ {
		if !(array[i] >= '0' && array[i] <= '9') {
			return false
		}
	}
	return true
}
