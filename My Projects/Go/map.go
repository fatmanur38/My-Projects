package piscine

func Map(f func(int) bool, a []int) []bool {
	r := make([]bool, len(a))
	for i, s := range a {
		r[i] = f(s)
	}
	return r
}
