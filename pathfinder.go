package main

import "sort"

// findPaths finds the optimal set of vertex-disjoint paths from start to end
// that minimises total turns for N ants.
// It uses DFS to find all simple paths, then exhaustive search to find
// the best vertex-disjoint combination.
func findPaths(farm *Farm) [][]string {
	allPaths := findAllSimplePaths(farm)

	if len(allPaths) == 0 {
		return nil
	}

	// Sort paths by length shortest first
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	return findBestDisjointSet(allPaths, farm.Ants)
}

// findAllSimplePaths uses DFS to find every simple path from start to end.
// A simple path visits each room at most once.
func findAllSimplePaths(farm *Farm) [][]string {
	result := [][]string{}
	visited := map[string]bool{farm.StartRoom: true}
	path := []string{farm.StartRoom}

	var dfs func(current string)
	dfs = func(current string) {
		if current == farm.EndRoom {
			p := make([]string, len(path))
			copy(p, path)
			result = append(result, p)
			return
		}

		// Sort neighbours for deterministic order
		neighbours := make([]string, len(farm.Tunnels[current]))
		copy(neighbours, farm.Tunnels[current])
		sort.Strings(neighbours)

		for _, next := range neighbours {
			if !visited[next] {
				visited[next] = true
				path = append(path, next)
				dfs(next)
				path = path[:len(path)-1]
				visited[next] = false
			}
		}
	}

	dfs(farm.StartRoom)
	return result
}

// findBestDisjointSet tries all combinations of vertex-disjoint paths
// and returns the subset that minimises total turns for numAnts ants.
// Two paths are vertex-disjoint if they share no intermediate rooms.
func findBestDisjointSet(paths [][]string, numAnts int) [][]string {
	bestTurns := -1
	bestSet := [][]string{}

	var try func(index int, chosen [][]string, usedRooms map[string]bool)
	try = func(index int, chosen [][]string, usedRooms map[string]bool) {
		if len(chosen) > 0 {
			turns := calculateTurns(chosen, numAnts)
			if bestTurns == -1 || turns < bestTurns {
				bestTurns = turns
				bestSet = make([][]string, len(chosen))
				copy(bestSet, chosen)
			}
		}

		for i := index; i < len(paths); i++ {
			path := paths[i]

			// Skip if this path shares any intermediate room with chosen paths
			conflict := false
			for _, room := range path[1 : len(path)-1] {
				if usedRooms[room] {
					conflict = true
					break
				}
			}
			if conflict {
				continue
			}

			// Mark intermediate rooms as used
			for _, room := range path[1 : len(path)-1] {
				usedRooms[room] = true
			}

			try(i+1, append(chosen, path), usedRooms)

			// Unmark intermediate rooms
			for _, room := range path[1 : len(path)-1] {
				delete(usedRooms, room)
			}
		}
	}

	try(0, [][]string{}, map[string]bool{})
	return bestSet
}

// calculateTurns simulates greedy ant assignment across paths and returns
// the actual number of turns needed to move all numAnts ants.
func calculateTurns(paths [][]string, numAnts int) int {
	if len(paths) == 0 {
		return 0
	}

	assigned := make([]int, len(paths))

	for ant := 0; ant < numAnts; ant++ {
		best := 0
		bestFinish := (len(paths[0]) - 1) + assigned[0] + 1
		for i := 1; i < len(paths); i++ {
			finish := (len(paths[i]) - 1) + assigned[i] + 1
			if finish < bestFinish {
				bestFinish = finish
				best = i
			}
		}
		assigned[best]++
	}

	maxTurns := 0
	for i, p := range paths {
		if assigned[i] == 0 {
			continue
		}
		finish := (len(p) - 1) + assigned[i]
		if finish > maxTurns {
			maxTurns = finish
		}
	}
	return maxTurns
}

// selectBestPaths returns the subset of paths that minimises total turns.
func selectBestPaths(allPaths [][]string, numAnts int) [][]string {
	bestTurns := -1
	bestCount := 1

	for i := 1; i <= len(allPaths); i++ {
		if i > numAnts {
			break
		}
		turns := calculateTurns(allPaths[:i], numAnts)
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			bestCount = i
		}
	}

	return allPaths[:bestCount]
}