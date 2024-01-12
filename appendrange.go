package piscine

func AppendRange(min, max int) []int {
	if min >= max {
		return nil
	}
	var sonuc []int
	for i := min; i < max; i++ {
		sonuc = append(sonuc, i)
	}
	return sonuc
}
