package piscine

func UltimateDivMod(a *int, b *int) {
	//*a = div
	//*b = mod
	//*a = *a / *b
	//div := *a / *b
	//mod := *a % *b

	var x int = *a
	// x := *a
	*a = *a / *b
	*b = x % *b
	// x := *a / *b
	//*b = *a % *b
}
