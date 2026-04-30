package main

// Move represents a single ant movement in one turn.
type Move struct {
	AntID int
	Room  string
}

// assignAnts distributes ants across the given paths.
// Ants are assigned greedily to whichever path would finish soonest.
func assignAnts(paths [][]string, numAnts int) [][]int {
	assignments := make([][]int, len(paths))

	for antID := 1; antID <= numAnts; antID++ {
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
	}

	return assignments
}

// turnsForPath calculates how many turns a path would take
// given it carries a certain number of ants sequentially.
func turnsForPath(path []string, numAnts int) int {
	steps := len(path) - 1
	return steps + numAnts - 1
}

// simulate generates all the turn-by-turn moves for the ants.
// Each ant at position antPos in its path queue departs on turn antPos
// (zero-indexed), moving one step per turn after that.
func simulate(paths [][]string, assignments [][]int) [][]Move {
	turns := [][]Move{}

	type antState struct {
		step int // current step index along its path (0 = still at start)
	}

	// Map antID -> its path and state
	antStates := map[int]*antState{}
	antPath := map[int][]string{}
	// antDepart[antID] = which turn (0-indexed) the ant is allowed to depart
	antDepart := map[int]int{}

	for pathIdx, ants := range assignments {
		for antPos, antID := range ants {
			antStates[antID] = &antState{step: 0}
			antPath[antID] = paths[pathIdx]
			// Each ant on the same path departs one turn after the previous
			antDepart[antID] = antPos
		}
	}

	turnIndex := 0
	for {
		turnMoves := []Move{}

		for _, ants := range assignments {
			for _, antID := range ants {
				state := antStates[antID]
				path := antPath[antID]

				// Ant can only move if:
				// 1. The current turn is >= its departure turn
				// 2. It has not yet reached the end of its path
				if turnIndex >= antDepart[antID] && state.step < len(path)-1 {
					state.step++
					turnMoves = append(turnMoves, Move{
						AntID: antID,
						Room:  path[state.step],
					})
				}
			}
		}

		if len(turnMoves) == 0 {
			break
		}

		turns = append(turns, turnMoves)
		turnIndex++
	}

	return turns
}