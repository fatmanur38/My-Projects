package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type farm struct {
	antAmount int
	rooms     map[string][2]int
	links     []string
	start     string
	end       string
	input     []string
}

type Node struct {
	Name        string
	Coordinates [2]int
}

type Edge struct {
	Start  *Node
	End    *Node
	Weight int
}

type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . example.txt")
		return
	}
	filename := os.Args[1]
	farminfo, err := readfile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if farminfo.antAmount <= 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	graph := createGraphFromRoomsAndLinks(farminfo.rooms, farminfo.links)
	if graph == nil {
		return
	}
	for _, line := range farminfo.input {
		fmt.Println(line)
	}
	fmt.Println()
	paths := findAllPaths(*graph, farminfo.start, farminfo.end)

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	var bestAntPaths [][]string
	minLength := -1

	for i := 0; i < len(paths); i++ {
		filterPaths := filterUniquePaths(paths, i)
		antPaths := moveAnts(filterPaths, farminfo.antAmount, farminfo.start, farminfo.end)

		if minLength == -1 || len(antPaths) < minLength {
			minLength = len(antPaths)
			bestAntPaths = antPaths
		}
	}

	printAntPaths(bestAntPaths)
}

func readfile(filename string) (farminfo farm, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return farminfo, fmt.Errorf("ERROR: cannot open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var filelines []string
	for scanner.Scan() {
		filelines = append(filelines, scanner.Text())
	}
	farminfo.input = filelines
	var links []string
	rooms := make(map[string][2]int)
	coordinates := make(map[[2]int]string)
	for i, line := range filelines {
		linecik := strings.Split(line, " ")
		if len(linecik) > 3 {
			return farminfo, fmt.Errorf("ERROR: invalid line format: %s", line)
		}
		if line == "##start" {
			start := strings.Fields(filelines[i+1])
			farminfo.start = start[0]
		} else if line == "##end" {
			end := strings.Fields(filelines[i+1])
			farminfo.end = end[0]
		} else if strings.Contains(line, "-") {
			links = append(links, line)
		} else if !(strings.Contains(line, "#")) && len(linecik) == 3 {
			x, err := strconv.Atoi(linecik[1])
			if err != nil {
				return farminfo, fmt.Errorf("ERROR: invalid coordinate: %s", linecik)
			}
			y, err := strconv.Atoi(linecik[2])
			if err != nil {
				return farminfo, fmt.Errorf("ERROR: invalid coordinate: %s", linecik)
			}
			coords := [2]int{x, y}
			if existingRoom, exists := coordinates[coords]; exists {
				return farminfo, fmt.Errorf("ERROR: duplicate coordinates: %s and %s have the same coordinates %v", linecik[0], existingRoom, coords)
			}
			coordinates[coords] = linecik[0]
			rooms[linecik[0]] = coords
		}

		if len(linecik) == 3 && RuneVarMi(linecik[0]) {
			return farminfo, fmt.Errorf("ERROR: invalid room name: %s", linecik)
		}
	}
	antamount, err := strconv.Atoi(filelines[0])
	if err != nil {
		return farminfo, fmt.Errorf("ERROR: invalid ant amount: %s", filelines[0])
	} else if antamount <= 0 {
		return farminfo, fmt.Errorf("ERROR: invalid data format")
	}
	farminfo.antAmount = antamount
	farminfo.rooms = rooms
	farminfo.links = links
	return farminfo, nil
}

func RuneVarMi(s string) bool {
	return s[0] == '#' || s[0] == 'L'
}

func pathContainsElement(path []string, elem string) bool {
	for _, item := range path {
		if item == elem {
			return true
		}
	}
	return false
}

func createGraphFromRoomsAndLinks(rooms map[string][2]int, links []string) *Graph {
	nodes := make(map[string]*Node)
	var edges []*Edge
	edgeSet := make(map[string]bool)

	for room, coords := range rooms {
		node := &Node{Name: room, Coordinates: coords}
		nodes[room] = node
	}

	for _, conn := range links {
		parts := strings.Split(conn, "-")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == parts[1] {
			fmt.Println("ERROR: invalid data format")
			return nil
		}
		startNode := nodes[parts[0]]
		endNode := nodes[parts[1]]
		if startNode != nil && endNode != nil {
			edgeKey := fmt.Sprintf("%s-%s", startNode.Name, endNode.Name)
			reverseEdgeKey := fmt.Sprintf("%s-%s", endNode.Name, startNode.Name)
			if edgeSet[edgeKey] || edgeSet[reverseEdgeKey] {
				fmt.Printf("ERROR: duplicate link found: %s\n", conn)
				return nil
			}
			edgeSet[edgeKey] = true
			edge1 := &Edge{Start: startNode, End: endNode, Weight: 1}
			edge2 := &Edge{Start: endNode, End: startNode, Weight: 1}
			edges = append(edges, edge1, edge2)
		}
	}

	graph := &Graph{}
	for _, node := range nodes {
		graph.Nodes = append(graph.Nodes, node)
	}
	graph.Edges = edges

	return graph
}

func findAllPaths(graph Graph, start string, end string) [][]string {
	visited := make(map[string]bool)
	paths := [][]string{}
	findPathsRecursive(graph, start, end, []string{start}, visited, &paths)
	return paths
}

func findPathsRecursive(graph Graph, current string, end string, path []string, visited map[string]bool, paths *[][]string) {
	visited[current] = true
	if current == end {
		newPath := make([]string, len(path))
		copy(newPath, path)
		*paths = append(*paths, newPath)
	} else {
		for _, edge := range graph.Edges {
			if edge.Start.Name == current && !visited[edge.End.Name] {
				findPathsRecursive(graph, edge.End.Name, end, append(path, edge.End.Name), visited, paths)
			}
		}
	}
	visited[current] = false
}

func filterUniquePaths(paths [][]string, index int) [][]string {
	firstPath := paths[index]
	midElements := make([]string, len(firstPath)-2)
	copy(midElements, firstPath[1:len(firstPath)-1])

	var uniquePaths [][]string
	uniquePaths = append(uniquePaths, firstPath)

	for i := 0; i < len(paths); i++ {
		if i == index {
			continue
		}
		path := paths[i]
		containsAny := false
		for _, elem := range midElements {
			if pathContainsElement(path, elem) {
				containsAny = true
				break
			}
		}

		if !containsAny {
			uniquePaths = append(uniquePaths, path)
			newMidElements := path[1 : len(path)-1]
			for _, newElem := range newMidElements {
				if !pathContainsElement(midElements, newElem) {
					midElements = append(midElements, newElem)
				}
			}
		}
	}
	return uniquePaths
}

func moveAnts(paths [][]string, antAmount int, startRoom string, endRoom string) [][]string {
	moves := [][]string{}
	antPositions := make([]int, antAmount)
	antPaths := make([][]string, antAmount)
	occupiedRooms := make(map[int]map[string]bool)
	occupiedEdges := make(map[int]map[string]bool)

	totalPathLength := 0
	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path)
		totalPathLength += len(path)
	}

	pathWeights := make([]float64, len(paths))
	antsPerPath := make([]int, len(paths))
	for i, length := range pathLengths {
		pathWeights[i] = float64(totalPathLength) / float64(length)
	}

	totalWeight := 0.0
	for _, weight := range pathWeights {
		totalWeight += weight
	}

	remainingAnts := antAmount
	for i := 0; i < len(paths); i++ {
		antsPerPath[i] = 1
		remainingAnts--
	}

	for remainingAnts > 0 {

		minIndex := -1
		minValue := int(^uint(0) >> 1)
		for i := 0; i < len(paths); i++ {
			value := (len(paths[i]) - 2) + antsPerPath[i]
			if value < minValue {
				minValue = value
				minIndex = i
			}
		}
		antsPerPath[minIndex]++
		remainingAnts--
	}

	antIndex := 0
	for i, ants := range antsPerPath {
		for j := 0; j < ants; j++ {
			antPaths[antIndex] = paths[i]
			antIndex++
		}
	}

	step := 0
	for {
		moveStep := []string{}
		allFinished := true

		if _, exists := occupiedRooms[step]; !exists {
			occupiedRooms[step] = make(map[string]bool)
		}
		if _, exists := occupiedEdges[step]; !exists {
			occupiedEdges[step] = make(map[string]bool)
		}

		for i := 0; i < antAmount; i++ {
			if antPositions[i] < len(antPaths[i])-1 {
				currentPosition := antPositions[i]
				nextPosition := currentPosition + 1
				currentRoom := antPaths[i][currentPosition]
				nextRoom := antPaths[i][nextPosition]
				edge := fmt.Sprintf("%s-%s", currentRoom, nextRoom)
				reverseEdge := fmt.Sprintf("%s-%s", nextRoom, currentRoom)

				if nextRoom != startRoom && nextRoom != endRoom && occupiedRooms[step][nextRoom] {
					continue
				}
				if occupiedEdges[step][edge] || occupiedEdges[step][reverseEdge] {
					continue
				}

				moveStep = append(moveStep, fmt.Sprintf("L%d-%s", i+1, nextRoom))
				occupiedRooms[step][nextRoom] = true
				occupiedEdges[step][edge] = true
				occupiedEdges[step][reverseEdge] = true
				antPositions[i]++
				allFinished = false
			}
		}

		if allFinished {
			break
		}

		moves = append(moves, moveStep)
		step++
	}

	return moves
}

func printAntPaths(antPaths [][]string) {
	for _, step := range antPaths {
		fmt.Println(strings.Join(step, " "))
	}
}
