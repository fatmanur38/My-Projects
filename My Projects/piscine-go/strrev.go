package piscine

func StrRev(s string) string {
	length := len(s)
	harf := []byte(s)
	// var x string
	for i := 0; i < length/2; i++ {
		// a := harf[i]
		// b := harf[length-i-1]
		// count := a
		// a = b
		// b = count
		harf[i], harf[length-i-1] = harf[length-i-1], harf[i]
		// x += string(harf[i])
	}
	return string(harf)
	/*package piscine

	func StrRev(s string) string {
		a := ""
		for i := len(s) - 1; i >= 0; i-- {
			a += string(s[i])
		}
		return a
	}*/
}
