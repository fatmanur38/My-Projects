package piscine

func ConcatParams(args []string) string {
	result := ""
	for i, c := range args {
		result += c
		if i != len(args)-1 {
			result += "\n"
		}
	}
	return result
}
