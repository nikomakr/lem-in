package main

// findPaths uses an Edmonds-Karp / BFS-based max flow approach to find
// the optimal set of non-overlapping paths from start to end.
// It returns the best subset of paths that minimises total turns for N ants.
func findPaths(farm *Farm) [][]string {
	// Build a capacity graph from the farm's list.
	// capacity[a][b] = 1 means the edge a->b is available.
	// capacity[a][b] = 0 means it is used.
	capacity := buildCapacityGraph(farm)

	// Find all paths using BFS (Edmonds-Karp)
	allPaths := [][]string{}
	for {
		path := bfs(farm.StartRoom, farm.EndRoom, capacity)
		if path == nil {
			break
		}
		allPaths = append(allPaths, path)

		// Update capacity along the found path:
		// - Reduce forward edge capacity to 0 (mark as used)
		// - Increase reverse edge capacity to 1 (allow undoing)
		for i := 0; i < len(path)-1; i++ {
			from := path[i]
			to := path[i+1]
			capacity[from][to]--
			capacity[to][from]++
		}
	}

	if len(allPaths) == 0 {
		return nil
	}

	// Select the optimal subset of paths that minimises turns for N ants
	return selectBestPaths(allPaths, farm.Ants)
}

// buildCapacityGraph creates a capacity map from the farm's tunnel list.
// Each undirected tunnel becomes two directed edges each with capacity 1.
// undirected edge a<->b becomes capacity[a][b] = 1 and capacity[b][a] = 1.
func buildCapacityGraph(farm *Farm) map[string]map[string]int {
	capacity := make(map[string]map[string]int)

	for room := range farm.Rooms {
		capacity[room] = make(map[string]int)
	}

	for from, neighbours := range farm.Tunnels {
		for _, to := range neighbours {
			capacity[from][to] = 1
		}
	}

	return capacity
}

// bfs performs a breadth-first search from start to end using available capacity.
// Returns the shortest path as a slice of room names, or nil if no path exists.
func bfs(start, end string, capacity map[string]map[string]int) []string {
	// Each element in the queue is a path (slice of room names)
	queue := [][]string{{start}}
	visited := map[string]bool{start: true}

	for len(queue) > 0 {
		// Dequeue the first path
		current := queue[0]
		queue = queue[1:]

		// The last room in the current path
		room := current[len(current)-1]

		if room == end {
			return current // found a path to the end
		}

		// Explore all neighbours with available capacity
		for neighbour, cap := range capacity[room] {
			if cap > 0 && !visited[neighbour] { 
				visited[neighbour] = true 
				// Build a new path by copying current and appending neighbour
				newPath := make([]string, len(current)+1)
				copy(newPath, current)
				newPath[len(current)] = neighbour
				queue = append(queue, newPath)
			}
		}
	}

	return nil // no path found
}

// calculateTurns returns the number of turns needed to move numAnts ants through a given set of paths.
// Formula: longest path length + (numAnts - numPaths)
// This assumes ants are distributed optimally across paths.
func calculateTurns(paths [][]string, numAnts int) int {
	if len(paths) == 0 {
		return 0
	}
	longest := 0
	for _, p := range paths {
		// Path length = number of steps = number of rooms - 1
		steps := len(p) - 1
		if steps > longest {
			longest = steps 
		}
	}
	return longest + (numAnts - len(paths)) // Extra turns for ants waiting to enter paths
}

// SelectBestPaths evaluates each prefix of allPaths (1 path, 2 paths, 3 paths...) and returns the subset that produces the fewest turns for numAnts ants.
// Adding more paths is only beneficial up to a point — beyond that, the longer path length outweighs the benefit of parallelism.
func selectBestPaths(allPaths [][]string, numAnts int) [][]string {
	bestTurns := -1 // bestTurns = -1 means we haven't found any valid configuration yet
	bestCount := 1 // bestCount = 1 means we start by considering just the first path

	for i := 1; i <= len(allPaths); i++ {
		// Only consider paths where we have at least one ant per path
		if i > numAnts {
			break
		}
		turns := calculateTurns(allPaths[:i], numAnts)
		// We want to minimize turns, so we update bestTurns and bestCount when we find a better configuration
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			bestCount = i
		}
	}

	return allPaths[:bestCount] // Return the best subset of paths that minimizes turns
}