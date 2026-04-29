package main

import "testing"

// --- assignAnts tests ---

func TestAssignAnts_SinglePath(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
	}
	assignments := assignAnts(paths, 3)
	if len(assignments[0]) != 3 {
		t.Errorf("expected 3 ants on path 0, got %d", len(assignments[0]))
	}
}

func TestAssignAnts_TwoEqualPaths(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	assignments := assignAnts(paths, 4)
	total := len(assignments[0]) + len(assignments[1])
	if total != 4 {
		t.Errorf("expected 4 ants total, got %d", total)
	}
}

func TestAssignAnts_MoreAntsOnShorterPath(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},                    // length 2 - shorter
		{"start", "b", "c", "d", "e", "f", "end"}, // length 6 - longer
	}
	assignments := assignAnts(paths, 6)
	// shorter path should carry more ants
	if len(assignments[0]) <= len(assignments[1]) {
		t.Errorf("expected more ants on shorter path, got %d vs %d",
			len(assignments[0]), len(assignments[1]))
	}
}

func TestAssignAnts_OneAnt(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	assignments := assignAnts(paths, 1)
	total := len(assignments[0]) + len(assignments[1])
	if total != 1 {
		t.Errorf("expected 1 ant total, got %d", total)
	}
}

// --- turnsForPath tests ---

func TestTurnsForPath_SingleAnt(t *testing.T) {
	path := []string{"start", "a", "b", "end"} // 3 steps
	turns := turnsForPath(path, 1)
	if turns != 3 {
		t.Errorf("expected 3 turns, got %d", turns)
	}
}

func TestTurnsForPath_MultipleAnts(t *testing.T) {
	path := []string{"start", "a", "end"} // 2 steps
	turns := turnsForPath(path, 3)
	// 2 + 3 - 1 = 4
	if turns != 4 {
		t.Errorf("expected 4 turns, got %d", turns)
	}
}

// --- simulate tests ---

func TestSimulate_SingleAntSinglePath(t *testing.T) {
	paths := [][]string{{"start", "a", "end"}}
	assignments := [][]int{{1}}

	turns := simulate(paths, assignments)

	if len(turns) == 0 {
		t.Fatal("expected at least one turn")
	}

	// Last move should land ant 1 at end
	lastTurn := turns[len(turns)-1]
	foundEnd := false
	for _, move := range lastTurn {
		if move.AntID == 1 && move.Room == "end" {
			foundEnd = true
		}
	}
	if !foundEnd {
		t.Error("expected ant 1 to reach 'end' in the last turn")
	}
}

func TestSimulate_AllAntsReachEnd(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	assignments := [][]int{{1, 2}, {3, 4}}

	turns := simulate(paths, assignments)

	// Count how many ants reached end across all turns
	reached := map[int]bool{}
	for _, turn := range turns {
		for _, move := range turn {
			if move.Room == "end" {
				reached[move.AntID] = true
			}
		}
	}
	if len(reached) != 4 {
		t.Errorf("expected all 4 ants to reach end, only %d did", len(reached))
	}
}

func TestSimulate_NoMovesOnEmptyAssignment(t *testing.T) {
	paths := [][]string{{"start", "end"}}
	assignments := [][]int{{}}

	turns := simulate(paths, assignments)
	if len(turns) != 0 {
		t.Errorf("expected 0 turns for empty assignment, got %d", len(turns))
	}
}