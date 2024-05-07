package piscine

func Capitalize(s string) string {
	sonuc := ""
	capitalizeNext := true

	for _, karakter := range s {
		if isAlphanumeric(karakter) {
			if capitalizeNext {
				sonuc += string(toUpper(karakter))
				capitalizeNext = false
			} else {
				sonuc += string(toLower(karakter))
			}
		} else {
			sonuc += string(karakter)
			capitalizeNext = true
		}
	}

	return sonuc
}

func isAlphanumeric(karakter rune) bool {
	return (karakter >= 'a' && karakter <= 'z') || (karakter >= 'A' && karakter <= 'Z') || (karakter >= '0' && karakter <= '9')
}

func toUpper(karakter rune) rune {
	if karakter >= 'a' && karakter <= 'z' {
		return karakter - 'a' + 'A'
	}
	return karakter
}

func toLower(karakter rune) rune {
	if karakter >= 'A' && karakter <= 'Z' {
		return karakter - 'A' + 'a'
	}
	return karakter
}
