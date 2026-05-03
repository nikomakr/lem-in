package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR: invalid data format, could not read file")
		os.Exit(1)
	}
	fileContent := string(fileBytes)

	farm, err := ParseFarm(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// DEBUG: print adjacency list
	fmt.Println("=== ADJACENCY LIST ===")
	keys := []string{}
	for k := range farm.Tunnels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s -> %v\n", k, farm.Tunnels[k])
	}

	paths := findPaths(farm)
	if paths == nil {
		fmt.Println("ERROR: invalid data format, no path found between start and end")
		os.Exit(1)
	}

	// DEBUG: print all found paths
	fmt.Println("=== PATHS FOUND ===")
	for i, p := range paths {
		fmt.Printf("  path %d: %v\n", i, p)
	}

	assignments := assignAnts(paths, farm.Ants)
	turns := simulate(paths, assignments)

	printOutput(fileContent, turns)
}