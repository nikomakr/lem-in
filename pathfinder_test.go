package main

import "testing"

// helper: builds a Farm from rooms and tunnels directly
func buildFarm(t *testing.T, ants int, start, end string, rooms []string, tunnels [][2]string) *Farm {
	t.Helper()
	farm := NewFarm()
	farm.Ants = ants
	farm.StartRoom = start
	farm.EndRoom = end

	for _, name := range rooms {
		r := &Room{Name: name}
		if name == start {
			r.IsStart = true
		}
		if name == end {
			r.IsEnd = true
		}
		farm.AddRoom(r)
	}
	for _, t2 := range tunnels {
		farm.AddTunnel(t2[0], t2[1])
	}
	return farm
}

// --- Single path ---

func TestFindPaths_SinglePath(t *testing.T) {
	farm := buildFarm(t, 3, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}, {"a", "end"}},
	)
	paths := findPaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected at least one path, got none")
	}
	if paths[0][0] != "start" || paths[0][len(paths[0])-1] != "end" {
		t.Errorf("expected path from start to end, got %v", paths[0])
	}
}

// --- Two parallel paths of equal length ---

func TestFindPaths_TwoParallelPaths(t *testing.T) {
	farm := buildFarm(t, 4, "start", "end",
		[]string{"start", "a", "b", "end"},
		[][2]string{
			{"start", "a"}, {"a", "end"},
			{"start", "b"}, {"b", "end"},
		},
	)
	paths := findPaths(farm)
	if len(paths) < 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
}

// --- No path exists ---

func TestFindPaths_NoPath(t *testing.T) {
	farm := buildFarm(t, 2, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}}, // no tunnel from a to end
	)
	paths := findPaths(farm)
	if paths != nil {
		t.Errorf("expected nil paths when no route exists, got %v", paths)
	}
}

// --- Start directly connected to end ---

func TestFindPaths_DirectConnection(t *testing.T) {
	farm := buildFarm(t, 1, "start", "end",
		[]string{"start", "end"},
		[][2]string{{"start", "end"}},
	)
	paths := findPaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected a direct path")
	}
	if len(paths[0]) != 2 {
		t.Errorf("expected path length 2, got %d", len(paths[0]))
	}
}

// --- calculateTurns ---

func TestCalculateTurns_SinglePath(t *testing.T) {
	paths := [][]string{{"start", "a", "b", "end"}} // length 3
	turns := calculateTurns(paths, 5)
	// 3 + (5 - 1) = 7
	if turns != 7 {
		t.Errorf("expected 7 turns, got %d", turns)
	}
}

func TestCalculateTurns_TwoPaths(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"}, // length 2
		{"start", "b", "end"}, // length 2
	}
	turns := calculateTurns(paths, 4)
	// 2 + (4 - 2) = 4
	if turns != 4 {
		t.Errorf("expected 4 turns, got %d", turns)
	}
}

func TestCalculateTurns_EmptyPaths(t *testing.T) {
	turns := calculateTurns([][]string{}, 5)
	if turns != 0 {
		t.Errorf("expected 0 turns for empty paths, got %d", turns)
	}
}

// --- selectBestPaths ---

func TestSelectBestPaths_PrefersShorterSubset(t *testing.T) {
	// Path 1: length 2, Path 2: length 2, Path 3: length 10
	// With 3 ants, adding path 3 makes things worse
	allPaths := [][]string{
		{"start", "a", "end"},                                                          // length 2
		{"start", "b", "end"},                                                          // length 2
		{"start", "c", "d", "e", "f", "g", "h", "i", "j", "k", "end"}, // length 10
	}
	best := selectBestPaths(allPaths, 3)
	if len(best) != 2 {
		t.Errorf("expected 2 paths selected, got %d", len(best))
	}
}

func TestSelectBestPaths_SingleAnt(t *testing.T) {
	allPaths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	best := selectBestPaths(allPaths, 1)
	// With 1 ant only 1 path makes sense
	if len(best) != 1 {
		t.Errorf("expected 1 path for 1 ant, got %d", len(best))
	}
}

// --- example00: 4 ants ---

func TestFindPaths_Example00(t *testing.T) {
	farm := buildFarm(t, 4, "0", "1",
		[]string{"0", "1", "2", "3"},
		[][2]string{
			{"0", "2"}, {"2", "3"}, {"3", "1"},
		},
	)
	paths := findPaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected paths for example00")
	}
	for _, p := range paths {
		if p[0] != "0" || p[len(p)-1] != "1" {
			t.Errorf("path does not go from start to end: %v", p)
		}
	}
}

// --- bfs directly ---

func TestBFS_FindsShortestPath(t *testing.T) {
	farm := buildFarm(t, 1, "start", "end",
		[]string{"start", "a", "b", "end"},
		[][2]string{
			{"start", "a"}, {"a", "end"},   // short path (length 2)
			{"start", "b"}, {"b", "end"},   // also length 2
			{"a", "b"},
		},
	)
	capacity := buildCapacityGraph(farm)
	path := bfs("start", "end", capacity)

	if path == nil {
		t.Fatal("expected BFS to find a path")
	}
	if path[0] != "start" || path[len(path)-1] != "end" {
		t.Errorf("expected path start->...->end, got %v", path)
	}
	if len(path) != 3 {
		t.Errorf("expected shortest path of length 3, got %d", len(path))
	}
}