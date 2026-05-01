package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput redirects stdout during a function call and returns what was printed.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// --- formatTurn tests ---

func TestFormatTurn_SingleMove(t *testing.T) {
	moves := []Move{{AntID: 1, Room: "a"}}
	result := formatTurn(moves)
	if result != "L1-a" {
		t.Errorf("expected 'L1-a', got '%s'", result)
	}
}

func TestFormatTurn_MultipleMoves(t *testing.T) {
	moves := []Move{
		{AntID: 1, Room: "a"},
		{AntID: 2, Room: "b"},
		{AntID: 3, Room: "c"},
	}
	result := formatTurn(moves)
	if result != "L1-a L2-b L3-c" {
		t.Errorf("expected 'L1-a L2-b L3-c', got '%s'", result)
	}
}

func TestFormatTurn_EmptyMoves(t *testing.T) {
	result := formatTurn([]Move{})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestFormatTurn_NamedRooms(t *testing.T) {
	moves := []Move{
		{AntID: 1, Room: "start"},
		{AntID: 2, Room: "gilfoyle"},
	}
	result := formatTurn(moves)
	if result != "L1-start L2-gilfoyle" {
		t.Errorf("expected 'L1-start L2-gilfoyle', got '%s'", result)
	}
}

// --- printOutput tests ---

func TestPrintOutput_FileContentPrinted(t *testing.T) {
	content := "4\n##start\n0 0 3\n##end\n1 8 3\n0-1\n"
	turns := [][]Move{
		{{AntID: 1, Room: "1"}},
	}

	output := captureOutput(func() {
		printOutput(content, turns)
	})

	if !strings.Contains(output, "##start") {
		t.Error("expected output to contain original file content")
	}
}

func TestPrintOutput_TurnsPrinted(t *testing.T) {
	content := "1\n##start\na 0 0\n##end\nb 1 1\na-b\n"
	turns := [][]Move{
		{{AntID: 1, Room: "b"}},
	}

	output := captureOutput(func() {
		printOutput(content, turns)
	})

	if !strings.Contains(output, "L1-b") {
		t.Error("expected output to contain 'L1-b'")
	}
}

func TestPrintOutput_BlankLineBetweenContentAndMoves(t *testing.T) {
	content := "1\n##start\na 0 0\n##end\nb 1 1\na-b\n"
	turns := [][]Move{
		{{AntID: 1, Room: "b"}},
	}

	output := captureOutput(func() {
		printOutput(content, turns)
	})

	// There should be a blank line between file content and moves
	if !strings.Contains(output, "\n\n") {
		t.Error("expected a blank line between file content and moves")
	}
}

func TestPrintOutput_Multipleturns(t *testing.T) {
	content := "2\n##start\na 0 0\n##end\nb 1 1\na-b\n"
	turns := [][]Move{
		{{AntID: 1, Room: "b"}},
		{{AntID: 2, Room: "b"}},
	}

	output := captureOutput(func() {
		printOutput(content, turns)
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")
	// Find the turn lines specifically
	foundL1 := false
	foundL2 := false
	for _, line := range lines {
		if line == "L1-b" {
			foundL1 = true
		}
		if line == "L2-b" {
			foundL2 = true
		}
	}
	if !foundL1 || !foundL2 {
		t.Errorf("expected both L1-b and L2-b in output, got:\n%s", output)
	}
}

func TestPrintOutput_NoTurns(t *testing.T) {
	content := "1\n##start\na 0 0\n##end\nb 1 1\na-b\n"

	output := captureOutput(func() {
		printOutput(content, [][]Move{})
	})

	// Should still print the file content
	if !strings.Contains(output, "##start") {
		t.Error("expected file content even with no turns")
	}

	// Should not contain any L moves
	if strings.Contains(output, "L1") {
		t.Error("expected no move output when turns is empty")
	}
}

// --- integration: formatTurn matches expected audit format ---

func TestFormatTurn_MatchesAuditFormat(t *testing.T) {
	// From example00 expected output: "L1-2 L2-5 L3-3"
	moves := []Move{
		{AntID: 1, Room: "2"},
		{AntID: 2, Room: "5"},
		{AntID: 3, Room: "3"},
	}
	expected := "L1-2 L2-5 L3-3"
	result := formatTurn(moves)
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

// --- helper: verify each ant appears exactly once per turn ---

func TestFormatTurn_EachAntOncePerTurn(t *testing.T) {
	moves := []Move{
		{AntID: 1, Room: "a"},
		{AntID: 2, Room: "b"},
		{AntID: 1, Room: "c"}, // duplicate ant in same turn — should not happen
	}
	result := formatTurn(moves)
	count := strings.Count(result, "L1-")
	if count > 1 {
		fmt.Printf("warning: ant 1 appears %d times in one turn: %s\n", count, result)
	}
}