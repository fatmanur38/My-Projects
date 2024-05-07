package piscine

func DescendAppendRange(max, min int) []int {
	if max <= min {
		return []int{}
	}
	var sonuc []int
	for i := max; i > min; i-- {
		sonuc = append(sonuc, i)
	}
	return sonuc
}
