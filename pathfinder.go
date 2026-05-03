package main

import "sort"

// findPaths uses Edmonds-Karp with node splitting to find
// the optimal set of non-overlapping paths from start to end.
func findPaths(farm *Farm) [][]string {
	initial := buildCapacityGraph(farm)
	capacity := cloneGraph(initial)

	foundAny := false
	for {
		path := bfsCapacity(farm.StartRoom+"_in", farm.EndRoom+"_out", capacity)
		if path == nil {
			break
		}
		foundAny = true
		for i := 0; i < len(path)-1; i++ {
			capacity[path[i]][path[i+1]]--
			capacity[path[i+1]][path[i]]++
		}
	}
	if !foundAny {
		return nil
	}

	realPaths := decomposeFlow(initial, capacity, farm.StartRoom, farm.EndRoom)
	if len(realPaths) == 0 {
		return nil
	}
	return selectBestPaths(realPaths, farm.Ants)
}

// cloneGraph returns a deep copy of a capacity graph.
func cloneGraph(g map[string]map[string]int) map[string]map[string]int {
	clone := make(map[string]map[string]int, len(g))
	for u, neighbors := range g {
		clone[u] = make(map[string]int, len(neighbors))
		for v, c := range neighbors {
			clone[u][v] = c
		}
	}
	return clone
}

// decomposeFlow extracts individual flow paths from the max-flow result.
// flow(u→v) = initial[u][v] - final[u][v]; we DFS through edges with flow > 0.
func decomposeFlow(initial, final map[string]map[string]int, startRoom, endRoom string) [][]string {
	flowGraph := make(map[string]map[string]int)
	for u, neighbors := range initial {
		for v, initCap := range neighbors {
			f := initCap - final[u][v]
			if f > 0 {
				if flowGraph[u] == nil {
					flowGraph[u] = make(map[string]int)
				}
				flowGraph[u][v] = f
			}
		}
	}

	startNode := startRoom + "_in"
	endNode := endRoom + "_out"
	paths := [][]string{}

	for {
		path := bfsFlow(startNode, endNode, flowGraph)
		if path == nil {
			break
		}
		paths = append(paths, toRealPath(path))
		for i := 0; i < len(path)-1; i++ {
			flowGraph[path[i]][path[i+1]]--
			if flowGraph[path[i]][path[i+1]] == 0 {
				delete(flowGraph[path[i]], path[i+1])
			}
		}
	}
	return paths
}

// bfsFlow finds a path in the flow graph (only forward/flow edges).
func bfsFlow(start, end string, flow map[string]map[string]int) []string {
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
		for neighbor, f := range flow[current] {
			if f > 0 {
				if _, visited := parent[neighbor]; !visited {
					parent[neighbor] = current
					queue = append(queue, neighbor)
				}
			}
		}
	}
	return nil
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

		for neighbour, c := range capacity[current] {
			if c > 0 {
				if _, visited := parent[neighbour]; !visited {
					parent[neighbour] = current
					queue = append(queue, neighbour)
				}
			}
		}
	}

	return nil
}

// toRealPath reconstructs real room names from a flow-decomposition path.
// Every room appears as X_in then X_out; we collect from _in nodes only.
func toRealPath(path []string) []string {
	result := []string{}
	for _, node := range path {
		if len(node) > 3 && node[len(node)-3:] == "_in" {
			result = append(result, node[:len(node)-3])
		}
	}
	return result
}

// calculateTurns returns the number of turns needed to move numAnts ants
// through a given set of paths.
func calculateTurns(paths [][]string, numAnts int) int {
	if len(paths) == 0 {
		return 0
	}
	longest := 0
	for _, p := range paths {
		steps := len(p) - 1
		if steps > longest {
			longest = steps
		}
	}
	return longest + (numAnts - len(paths))
}

// selectBestPaths returns the subset of paths that minimises total turns.
// Paths are sorted shortest-first before selection so prefix-comparison is valid.
func selectBestPaths(allPaths [][]string, numAnts int) [][]string {
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

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