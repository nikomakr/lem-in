package main

import (
	"fmt"
	"strings"
)

// formatTurn converts a slice of moves into the "Lx-y Lx-y" string format.
func formatTurn(moves []Move) string {
	parts := make([]string, len(moves))
	for i, move := range moves {
		parts[i] = fmt.Sprintf("L%d-%s", move.AntID, move.Room)
	}
	return strings.Join(parts, " ")
}

// printOutput prints the original file content followed by each turn.
func printOutput(fileContent string, turns [][]Move) {
	// Print the original file content
	fmt.Print(fileContent)

	// Print a blank line between file content and moves
	fmt.Println()

	// Print each turn on its own line
	for _, turn := range turns {
		fmt.Println(formatTurn(turn))
	}
}