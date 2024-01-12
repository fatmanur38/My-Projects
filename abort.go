package piscine

func Abort(a, b, c, d, e int) int {
	table := []int{a, b, c, d, e}
	SortIntegerTablex(table)
	return table[2]
}

func SortIntegerTablex(table []int) {
	for i := 0; i < len(table); i++ {
		for j := i + 1; j < len(table); j++ {
			if table[i] > table[j] {
				gecici := table[j]
				table[j] = table[i]
				table[i] = gecici
			}
		}
	}
}
