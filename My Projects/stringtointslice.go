package piscine

func StringToIntSlice(str string) []int {
	sonuc := []int{}
	if str == "" {
		return nil
	}
	for _, karakter := range str {
		sonuc = append(sonuc, int(karakter))
	}
	return sonuc
}
