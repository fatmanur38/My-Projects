package piscine

func JumpOver(str string) string {
	sonuc := ""
	for i := 2; i < len(str); i += 3 {
		sonuc += string(str[i])
	}
	sonuc += "\n"
	return sonuc
}
