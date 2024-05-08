package piscine

func IsAlpha(s string) bool {
	array := []rune(s)
	for i := 0; i < len(s); i++ {
		if !((array[i] >= 'a' && array[i] <= 'z') || (array[i] >= 'A' && array[i] <= 'Z') || (array[i] >= '0' && array[i] <= '9')) {
			return false
		}
	}
	return true
}
