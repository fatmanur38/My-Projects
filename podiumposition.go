package piscine

func PodiumPosition(podium [][]string) [][]string {
	for i := 0; i < len(podium); i++ {
		for j := 0; j < len(podium[i]); j++ {
			if j == 0 {
				podium[i][j] += "st"
			} else if j == 1 {
				podium[i][j] += "nd"
			} else if j == 2 {
				podium[i][j] += "rd"
			} else {
				podium[i][j] += "th"
			}
		}
	}
	return podium
}
