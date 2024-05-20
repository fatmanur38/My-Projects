package piscine

func Rot14(s string) string {
	result := ""
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			base := 'a'
			if char >= 'A' && char <= 'Z' {
				base = 'A'
			}
			result += string((char-base+14)%26 + base)
		} else {
			result += string(char)
		}
	}
	return result
}
