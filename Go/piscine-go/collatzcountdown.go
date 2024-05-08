package piscine

func CollatzCountdown(start int) int {
	if start <= 0 {
		return -1
	}
	if start == 1 {
		return 0
	}
	if start%2 == 0 {
		return 1 + CollatzCountdown(start/2)
	}
	return 1 + CollatzCountdown(start*3+1)
}
