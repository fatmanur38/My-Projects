package piscine

func Compact(ptr *[]string) int {
	slice := *ptr
	size := 0

	for _, v := range slice {
		if v != "" {
			slice[size] = v
			size++
		}
	}

	*ptr = slice[:size]
	return size
}

/*package piscine

func Compact(ptr *[]string) int {
	count := 0
	x := 0
	for i := 0; i < len(*ptr); i++ {
		if (*ptr)[i] != "" {
			count++
		}
	}
	str := make([]string, count)
	for i := 0; i < len(*ptr); i++ {
		if (*ptr)[i] != "" {
			str[x] = (*ptr)[i]
			x++
		}
	}
	(*ptr) = str[:]
	return count
}*/
