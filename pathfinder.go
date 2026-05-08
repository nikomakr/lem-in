package main

import "sort"

// findPaths uses Edmonds-Karp with node splitting to find
// the optimal set of non-overlapping paths from start to end.
func findPaths(farm *Farm) [][]string {
	capacity := buildCapacityGraph(farm)

	rawPaths := [][]string{}
	for {
		path := bfsCapacity(farm.StartRoom+"_in", farm.EndRoom+"_out", capacity)
		if path == nil {
			break
		}
		rawPaths = append(rawPaths, path)

		for i := 0; i < len(path)-1; i++ {
			from := path[i]
			to := path[i+1]
			capacity[from][to]--
			capacity[to][from]++
		}
	}

	if len(rawPaths) == 0 {
		return nil
	}

	realPaths := [][]string{}
	for _, p := range rawPaths {
		realPaths = append(realPaths, toRealPath(p))
	}

	return selectBestPaths(realPaths, farm.Ants)
}

// buildCapacityGraph builds a residual capacity graph using node splitting.
// Each room X becomes X_in and X_out with an internal edge.
// Start and end get capacity equal to their tunnel count.
// All other rooms get internal capacity 1.
func buildCapacityGraph(farm *Farm) map[string]map[string]int {
	cap := make(map[string]map[string]int)

	ensure := func(node string) {
		if cap[node] == nil {
			cap[node] = make(map[string]int)
		}
	}

	for room := range farm.Rooms {
		in := room + "_in"
		out := room + "_out"
		ensure(in)
		ensure(out)

		if room == farm.StartRoom || room == farm.EndRoom {
			cap[in][out] = len(farm.Tunnels[room])
		} else {
			cap[in][out] = 1
		}
		cap[out][in] = 0
	}

	for from, neighbours := range farm.Tunnels {
		for _, to := range neighbours {
			fromOut := from + "_out"
			toIn := to + "_in"
			ensure(fromOut)
			ensure(toIn)
			cap[fromOut][toIn] = 1
			if _, exists := cap[toIn][fromOut]; !exists {
				cap[toIn][fromOut] = 0
			}
		}
	}

	return cap
}

// bfsCapacity finds the shortest augmenting path using parent tracking BFS.
// Neighbours are sorted before exploration to ensure deterministic results
// across runs, since Go map iteration order is randomised.
func bfsCapacity(start, end string, capacity map[string]map[string]int) []string {
	parent := map[string]string{start: ""}
	queue := []string{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			path := []string{}
			for node := end; node != ""; node = parent[node] {
				path = append([]string{node}, path...)
			}
			return path
		}

		// Sort neighbours for deterministic BFS order
		neighbours := make([]string, 0, len(capacity[current]))
		for neighbour, c := range capacity[current] {
			if c > 0 {
				neighbours = append(neighbours, neighbour)
			}
		}
		sort.Strings(neighbours)

		for _, neighbour := range neighbours {
			if _, visited := parent[neighbour]; !visited {
				parent[neighbour] = current
				queue = append(queue, neighbour)
			}
		}
	}

	return nil
}

// toRealPath reconstructs real room names from a split-node path.
// Collects _in nodes for forward traversal and appends the final
// _out node only if not already present.
func toRealPath(path []string) []string {
	result := []string{}
	seen := map[string]bool{}

	for _, node := range path {
		if len(node) > 3 && node[len(node)-3:] == "_in" {
			name := node[:len(node)-3]
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}

	last := path[len(path)-1]
	if len(last) > 4 && last[len(last)-4:] == "_out" {
		name := last[:len(last)-4]
		if !seen[name] {
			result = append(result, name)
		}
	}

	return result
}

// calculateTurns simulates greedy ant assignment across paths and returns
// the actual number of turns needed to move all numAnts ants.
// Each ant is assigned to whichever path finishes it soonest.
// This correctly handles unequal path lengths unlike the simplified formula.
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