package main

import (
	"os"
)

func printStr(s string) {
	os.Stdout.Write([]byte(s))
}

func main() {
	// Komut satırı argümanlarını kontrol et
	args := os.Args[1:]
	if len(args) == 0 {
		printStr("File name missing\n")
		return
	} else if len(args) > 1 {
		printStr("Too many arguments\n")
		return
	}

	// Dosya adını al
	fileName := args[0]

	// Dosyayı aç
	file, err := os.Open(fileName)
	if err != nil {
		printStr("Error opening file: " + err.Error() + "\n")
		return
	}
	defer file.Close()

	// Dosya içeriğini oku
	content := make([]byte, 1024) // Örnek: Maksimum 1024 byte okuma
	n, err := file.Read(content)
	if err != nil {
		printStr("Error reading file: " + err.Error() + "\n")
		return
	}

	// Dosya içeriğini yazdır
	os.Stdout.Write(content[:n])
}
