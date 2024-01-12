package main

import "os"

func atoi(s string) (int, bool) {
	sum := 0
	minus := false
	for i, c := range s {
		if c == '-' && i == 0 {
			minus = true
		} else if c >= 48 && c <= 57 {
			l := toInt(c)
			sum = sum*10 + l
		} else {
			return 0, false
		}
	}
	if minus {
		sum = -sum
	}
	return sum, true
}

func toInt(a rune) int {
	s := 0
	for i := '0'; i < a; i++ {
		s++
	}
	return s
}

func itoa(n int) (s string) {
	for i := n; i != 0; i = i / 10 {
		zz := i % 10
		if zz < 0 {
			zz = -zz
		}
		s = string(rune(zz+48)) + s
	}
	if n < 0 {
		s = "-" + s
	}
	if n == 0 {
		s = "0"
	}
	return
}

func main() {
	l := 0
	for range os.Args {
		l++
	}
	if l != 4 {
		return
	}
	a, ok := atoi(os.Args[1])
	if !ok {
		return
	}
	b, ok := atoi(os.Args[3])
	if !ok {
		return
	}
	switch os.Args[2] {
	case "+":
		result := a + b
		if (result > a) == (b > 0) {
			os.Stdout.WriteString(itoa(result) + "\n")
		}
	case "-":
		result := a - b
		if (result < a) == (b > 0) {
			os.Stdout.WriteString(itoa(result) + "\n")
		}
	case "*":
		result := a * b
		if a == 0 || (result/a == b) {
			os.Stdout.WriteString(itoa(result) + "\n")
		}
	case "/":
		if b == 0 {
			os.Stdout.WriteString("No division by 0\n")
		} else {
			os.Stdout.WriteString(itoa(a/b) + "\n")
		}
	case "%":
		if b == 0 {
			os.Stdout.WriteString("No modulo by 0\n")
		} else {
			os.Stdout.WriteString(itoa(a%b) + "\n")
		}
	}
}
