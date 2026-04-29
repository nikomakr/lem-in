package main

// Move represents a single ant movement in one turn.
type Move struct {
	AntID int
	Room  string
}

// AssignAnts distributes ants across the given paths.
// Returns a slice where each index i contains the ant number assigned to path i.
// Ants are assigned one per path in round-robin order.
func assignAnts(paths [][]string, numAnts int) [][]int {
	assignments := make([][]int, len(paths))

	// Assign ants to paths based on which path will finish soonest
	antID := 1
	for antID <= numAnts {
		// Pick the path that would finish the earliest if we added one more ant
		bestPath := 0
		bestTurns := turnsForPath(paths[0], len(assignments[0])+1)

		for i := 1; i < len(paths); i++ {
			t := turnsForPath(paths[i], len(assignments[i])+1)
			if t < bestTurns {
				bestTurns = t
				bestPath = i
			}
		}
		assignments[bestPath] = append(assignments[bestPath], antID)
		antID++
	}

	return assignments
}

// turnsForPath calculates how many turns a path would take given it carries a certain number of ants sequentially.
func turnsForPath(path []string, numAnts int) int {
	steps := len(path) - 1
	return steps + numAnts - 1
}

// Simulate generates all the turn-by-turn moves for the ants.
// Returns a slice of turns, each turn being a slice of Move.
func simulate(paths [][]string, assignments [][]int) [][]Move {
	turns := [][]Move{}

	// Track each ant's current position along its path
	// antProgress[antID] = current step index along its path
	type antState struct {
		pathIndex int
		step      int
	}

	// Build a map of antID -> its path and current step
	antStates := map[int]*antState{}
	antPath := map[int][]string{}

	for pathIdx, ants := range assignments {
		for _, antID := range ants {
			antStates[antID] = &antState{pathIndex: pathIdx, step: 0}
			antPath[antID] = paths[pathIdx]
		}
	}

	// Run simulation until all ants have reached the end
	for {
		turnMoves := []Move{}

		for pathIdx, ants := range assignments {
			for antPos, antID := range ants {
				state := antStates[antID]
				path := antPath[antID]

				// An ant can only move if the next room is free
				// Ants on the same path are spaced one step apart
				// Ant at position antPos can move on turn: antPos + step
				nextStep := state.step + 1

				// Check if the ant can move to the next step
				if nextStep <= antPos+len(turns) && nextStep < len(path) {
					state.step = nextStep
					turnMoves = append(turnMoves, Move{
						AntID: antID,
						Room:  path[nextStep],
					})
				}
				_ = pathIdx
			}
		}

		if len(turnMoves) == 0 {
			break
		}

		turns = append(turns, turnMoves)
	}

	return turns
}