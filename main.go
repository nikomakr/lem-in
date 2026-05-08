package main

import (
	"fmt"
	"os"
)

func main() {
	// 1. Read filename from command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	// 2. Parse the farm — file content is returned alongside the farm
	farm, fileContent, err := ParseFarm(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Find optimal paths
	paths := findPaths(farm)
	if paths == nil {
		fmt.Println("ERROR: invalid data format, no path found between start and end")
		os.Exit(1)
	}

	// 4. Assign ants to paths and simulate movement
	assignments := assignAnts(paths, farm.Ants)
	turns := simulate(paths, assignments)

	// 5. Print output
	printOutput(fileContent, turns)
}