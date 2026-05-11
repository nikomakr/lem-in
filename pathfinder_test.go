package main

import "testing"

func buildTestFarm(t *testing.T, ants int, start, end string, rooms []string, tunnels [][2]string) *Farm {
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
	for _, tun := range tunnels {
		farm.AddTunnel(tun[0], tun[1])
	}
	return farm
}

// --- findAllSimplePaths ---

func TestFindAllSimplePaths_SinglePath(t *testing.T) {
	farm := buildTestFarm(t, 1, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}, {"a", "end"}},
	)
	paths := findAllSimplePaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected at least one path")
	}
}

func TestFindAllSimplePaths_TwoParallelPaths(t *testing.T) {
	farm := buildTestFarm(t, 2, "start", "end",
		[]string{"start", "a", "b", "end"},
		[][2]string{
			{"start", "a"}, {"a", "end"},
			{"start", "b"}, {"b", "end"},
		},
	)
	paths := findAllSimplePaths(farm)
	if len(paths) < 2 {
		t.Errorf("expected at least 2 paths, got %d", len(paths))
	}
}

func TestFindAllSimplePaths_NoPath(t *testing.T) {
	farm := buildTestFarm(t, 1, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}},
	)
	paths := findAllSimplePaths(farm)
	if len(paths) != 0 {
		t.Errorf("expected no paths, got %d", len(paths))
	}
}

// --- findBestDisjointSet ---

func TestFindBestDisjointSet_TwoDisjointPaths(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	best := findBestDisjointSet(paths, 4)
	if len(best) != 2 {
		t.Errorf("expected 2 disjoint paths, got %d", len(best))
	}
}

func TestFindBestDisjointSet_ConflictingPaths(t *testing.T) {
	paths := [][]string{
		{"start", "a", "b", "end"},
		{"start", "a", "c", "end"},
		{"start", "d", "end"},
	}
	best := findBestDisjointSet(paths, 4)
	for i := 0; i < len(best)-1; i++ {
		for j := i + 1; j < len(best); j++ {
			roomsI := map[string]bool{}
			for _, r := range best[i][1 : len(best[i])-1] {
				roomsI[r] = true
			}
			for _, r := range best[j][1 : len(best[j])-1] {
				if roomsI[r] {
					t.Errorf("paths %d and %d share room %s", i, j, r)
				}
			}
		}
	}
}

// --- calculateTurns ---

func TestCalculateTurns_SinglePath(t *testing.T) {
	paths := [][]string{{"start", "a", "b", "end"}} // 3 steps
	turns := calculateTurns(paths, 5)
	// finish = steps + numAnts = 3 + 5 = 8
	if turns != 8 {
		t.Errorf("expected 8 turns, got %d", turns)
	}
}

func TestCalculateTurns_TwoEqualPaths(t *testing.T) {
	paths := [][]string{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	turns := calculateTurns(paths, 4)
	if turns != 4 {
		t.Errorf("expected 4 turns, got %d", turns)
	}
}

func TestCalculateTurns_EmptyPaths(t *testing.T) {
	turns := calculateTurns([][]string{}, 5)
	if turns != 0 {
		t.Errorf("expected 0 turns, got %d", turns)
	}
}

// --- findPaths ---

func TestFindPaths_SinglePath(t *testing.T) {
	farm := buildTestFarm(t, 3, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}, {"a", "end"}},
	)
	paths := findPaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected at least one path")
	}
	if paths[0][0] != "start" || paths[0][len(paths[0])-1] != "end" {
		t.Errorf("path does not go from start to end: %v", paths[0])
	}
}

func TestFindPaths_NoPath(t *testing.T) {
	farm := buildTestFarm(t, 2, "start", "end",
		[]string{"start", "a", "end"},
		[][2]string{{"start", "a"}},
	)
	paths := findPaths(farm)
	if paths != nil {
		t.Errorf("expected nil, got %v", paths)
	}
}

func TestFindPaths_TwoParallelPaths(t *testing.T) {
	farm := buildTestFarm(t, 4, "start", "end",
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

func TestFindPaths_Example00(t *testing.T) {
	farm := buildTestFarm(t, 4, "0", "1",
		[]string{"0", "1", "2", "3"},
		[][2]string{{"0", "2"}, {"2", "3"}, {"3", "1"}},
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

func TestFindPaths_ThreeDisjointPaths(t *testing.T) {
	farm := buildTestFarm(t, 10, "start", "end",
		[]string{"start", "end", "a", "b", "c", "d", "e", "f"},
		[][2]string{
			{"start", "a"}, {"a", "end"},
			{"start", "b"}, {"b", "end"},
			{"start", "c"}, {"c", "d"}, {"d", "e"}, {"e", "f"}, {"f", "end"},
		},
	)
	paths := findPaths(farm)
	if len(paths) == 0 {
		t.Fatal("expected paths")
	}
	usedRooms := map[string]bool{}
	for _, path := range paths {
		for _, room := range path[1 : len(path)-1] {
			if usedRooms[room] {
				t.Errorf("room %s used in multiple paths", room)
			}
			usedRooms[room] = true
		}
	}
}