package piscine

func IsPrintable(s string) bool {
	for _, karakter := range s {
		if karakter < 32 || karakter > 126 {
			return false
		}
	}
	return true
}
