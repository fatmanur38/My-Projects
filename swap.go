package piscine

func Swap(a *int, b *int) {
	gecici := *a
	*a = *b
	*b = gecici
}
