package piscine

func ForEach(f func(int), a []int) {
	for i := range a {
		f(a[i])
	}
}

/*
	for i:=0;i<len(a);i++{
		f(a[i])
	}


*/
