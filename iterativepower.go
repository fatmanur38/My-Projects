package piscine

func IterativePower(nb int, power int) int {
	if power < 0 {
		return 0
	}
	sonuc := 1
	for i := 1; i <= power; i++ {
		sonuc *= nb
	}
	return sonuc
}
