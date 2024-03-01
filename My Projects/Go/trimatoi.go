package piscine

func TrimAtoi(s string) int {
	result := 0
	sign := 1
	foundDigit := false

	for _, char := range s {
		if char == '-' && !foundDigit {
			sign = -1
		} else if char >= '0' && char <= '9' {
			digit := int(char - '0')
			result = result*10 + digit
			foundDigit = true
		}
	}

	return result * sign
}
