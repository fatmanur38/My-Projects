package piscine

func ReverseMenuIndex(menu []string) []string {
	result := []string{}
	for i := len(menu) - 1; i >= 0; i-- {
		result = append(result, menu[i])
	}
	return result
}
