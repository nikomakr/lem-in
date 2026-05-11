package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	farm, fileContent, err := ParseFarm(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paths := findPaths(farm)
	if paths == nil {
		fmt.Println("ERROR: invalid data format, no path found between start and end")
		os.Exit(1)
	}

	assignments := assignAnts(paths, farm.Ants)

	// // DEBUG
	// for i, p := range paths {
	// 	fmt.Printf("DEBUG path %d (len %d): %v\n", i, len(p), p)
	// }
	// for i, a := range assignments {
	// 	fmt.Printf("DEBUG assignment %d: %v\n", i, a)
	// }

	turns := simulate(paths, assignments)
	printOutput(fileContent, turns)
}