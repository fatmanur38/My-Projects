package piscine

func SplitWhiteSpaces(s string) []string {
	var words []string
	hiclik := ""
	for _, karakter := range s {
		if karakter == ' ' || karakter == '\t' || karakter == '\n' {
			if hiclik != "" {
				words = append(words, hiclik)
				hiclik = ""
			}
		} else {
			hiclik += string(karakter)
		}
	}
	if hiclik != "" {
		words = append(words, hiclik)
	}
	return words
}
