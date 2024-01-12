package piscine

func ToUpper(s string) string {
	sonuc := ""
	for _, karakter := range s {
		if karakter >= 'a' && karakter <= 'z' {
			sonuc += string(karakter - 32)
		} else {
			sonuc += string(karakter)
		}
	}
	return sonuc
}
