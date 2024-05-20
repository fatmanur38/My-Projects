package piscine

func BasicJoin(elems []string) string {
	// Tüm stringleri birleştirmek için bir döngü kullanalım.
	var result string
	for _, elem := range elems {
		result += elem
	}
	return result
}
