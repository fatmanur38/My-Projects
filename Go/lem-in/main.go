package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type farm struct {
	antAmount int
	rooms     []string
	links     []string
	start     int
	end       int
}

func main() {
	filename := os.Args[1]
	fmt.Println(readfile(filename))
}

func readfile(filename string) (farminfo farm) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var filelines []string
	startin := 0
	endin := 0

	for scanner.Scan() {
		filelines = append(filelines, scanner.Text())
	}
	antamountstr := ""
	var links []string
	var rooms []string
	for i, line := range filelines {
		linecik := strings.Split(line, " ")
		if line == "##start" {
			startin = i
		} else if line == "##end" {
			endin = i
		} else if strings.Contains(line, "-") {
			links = append(links, line)
		} else if !(strings.Contains(line, "#")) && len(linecik) == 3 {
			rooms = append(rooms, line)
		} else {
			antamountstr = line
		}
	}
	antamount, err := strconv.Atoi(antamountstr)
	if err != nil {
		fmt.Println("error")
	}

	farminfo.antAmount = antamount
	farminfo.start = startin
	farminfo.end = endin
	farminfo.rooms = rooms
	farminfo.links = links
	return farminfo
}

func atoi(s string) int {
	var sonuc int

	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			sonuc = 10*sonuc + int(ch-'0')
		}
	}
	return sonuc
}
