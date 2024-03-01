package piscine

func RockAndRoll(n int) string {
	if n < 0 {
		return "error: number is negative\n"
	}

	if n%2 == 0 && n%3 == 0 {
		return "rock and roll\n"
	} else if n%2 == 0 {
		return "rock\n"
	} else if n%3 == 0 {
		return "roll\n"
	} else {
		return "error: non divisible\n"
	}
}
