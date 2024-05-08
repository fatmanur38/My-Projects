package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", AsciiToWeb)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func AsciiToWeb(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./assets/index.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		str := r.FormValue("text")
		for _, char := range str {
			for _, turkce := range "çöşıüğİÇŞÜĞÖ" {
				if char == turkce {
					http.Error(w, "Turkish characters are not allowed", http.StatusBadRequest)
					return
				}
			}
		}
		banner := r.FormValue("banner")
		file, err := os.Open(banner + ".txt") // to open file
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		str1 := strings.Split(str, "\n")
		var res []string
		for _, wordgroup := range str1 {
			res = StrWord(wordgroup, file)
			for _, line := range res {
				fmt.Fprintf(w, "%s\n", line)
			}
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func StrWord(word string, file *os.File) []string {
	lines := make([]string, 8)
	for _, char := range word {
		for i := 0; i <= 7; i++ { // the string (lines) in our code needs 8 lines
			file.Seek(0, 0)                   // to reset the file, in order to read again
			start := (int(char)-32)*9 + 2 + i // to match the lines in "standard.txt" with which the char we took in our word
			reader := bufio.NewScanner(file)
			satir := 1
			for reader.Scan() {
				if satir == start {
					lines[i] += reader.Text() // to add the line to lines
				}
				satir++
			}
			if err := reader.Err(); err != nil {
				fmt.Println("Dosya okunurken bir hata oluştu:", err)
				return nil
			}
		}
	}
	return lines
}
