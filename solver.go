package main

// Move represents a single ant movement in one turn.
type Move struct {
	AntID int
	Room  string
}

// assignAnts distributes ants across the given paths using greedy assignment.
// Each ant is assigned to whichever path finishes it soonest.
// This matches the greedy logic in calculateTurns in pathfinder.go.
func assignAnts(paths [][]string, numAnts int) [][]int {
	assignments := make([][]int, len(paths))
	assigned := make([]int, len(paths))

	for antID := 1; antID <= numAnts; antID++ {
		best := 0
		bestFinish := (len(paths[0]) - 1) + assigned[0] + 1

		for i := 1; i < len(paths); i++ {
			finish := (len(paths[i]) - 1) + assigned[i] + 1
			if finish < bestFinish {
				bestFinish = finish
				best = i
			}
		}
		assignments[best] = append(assignments[best], antID)
		assigned[best]++
	}

	return assignments
}

// turnsForPath calculates how many turns a path takes for a given number of ants.
func turnsForPath(path []string, numAnts int) int {
	steps := len(path) - 1
	return steps + numAnts - 1
}

// simulate generates all the turn-by-turn moves for the ants.
// Each ant at position antPos in its path queue departs on turn antPos.
// A room occupancy map ensures no two ants occupy the same room per turn.
func simulate(paths [][]string, assignments [][]int) [][]Move {
	turns := [][]Move{}

	type antState struct {
		step int
	}

	antStates := map[int]*antState{}
	antPath := map[int][]string{}
	antDepart := map[int]int{}

	for pathIdx, ants := range assignments {
		for antPos, antID := range ants {
			antStates[antID] = &antState{step: 0}
			antPath[antID] = paths[pathIdx]
			antDepart[antID] = antPos
		}
	}

	// Collect all ant IDs in order for deterministic iteration
	allAnts := []int{}
	for _, ants := range assignments {
		allAnts = append(allAnts, ants...)
	}

	turnIndex := 0
	for {
		turnMoves := []Move{}

		// Track rooms occupied this turn to enforce one-ant-per-room
		// End room is exempt — multiple ants can arrive there
		occupied := map[string]bool{}

		for _, antID := range allAnts {
			state := antStates[antID]
			path := antPath[antID]

			if turnIndex >= antDepart[antID] && state.step < len(path)-1 {
				nextRoom := path[state.step+1]
				isEnd := nextRoom == path[len(path)-1]

				// Allow move if room is end (exempt) or not yet occupied
				if isEnd || !occupied[nextRoom] {
					if !isEnd {
						occupied[nextRoom] = true
					}
					state.step++
					turnMoves = append(turnMoves, Move{
						AntID: antID,
						Room:  nextRoom,
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