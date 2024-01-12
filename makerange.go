package piscine

func MakeRange(min, max int) []int {
	if min >= max {
		return nil
	}
	kacsayi := max - min
	sonuc := make([]int, kacsayi)
	for i := 0; i < kacsayi; i++ {
		sonuc[i] = min + i
	}
	return sonuc
}
