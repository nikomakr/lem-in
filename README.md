# lem-in

A Go program that solves an ant colony pathfinding problem. Parses a graph of rooms and tunnels, finds optimal non-overlapping paths using BFS-based max flow, and moves N ants from `##start` to `##end` in the fewest turns possible.

---

## Table of Contents

- [Overview](#overview)
- [Algorithm](#algorithm)
- [Project Structure](#project-structure)
- [Input Format](#input-format)
- [Output Format](#output-format)
- [Usage](#usage)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Testing](#testing)

---

## Overview

`lem-in` simulates a digital ant farm. Given a colony of rooms connected by tunnels, the program finds the most efficient way to move **N ants** from the `##start` room to the `##end` room in as few turns as possible.

Rules:
- Each room holds at most **one ant at a time** (except `##start` and `##end`)
- Each tunnel can only be used **once per turn**
- All ants begin in `##start` and must reach `##end`

---

## Algorithm

This project uses an **Edmonds-Karp / BFS-based Max Flow** approach rather than a simple shortest-path algorithm like Dijkstra. Here is why:

### Why Not Dijkstra?

Dijkstra finds a single shortest path between two nodes. For `lem-in`, that is not enough — we need to move **N ants simultaneously** across **multiple non-overlapping paths**, minimising the total number of turns. Dijkstra cannot model traffic, congestion, or concurrent ant movement.

### The Edmonds-Karp Approach

Edmonds-Karp is an implementation of the Ford-Fulkerson max flow algorithm that uses **BFS to find augmenting paths**. In the context of `lem-in`:

1. **BFS** finds the shortest available path from `##start` to `##end`
2. The path is recorded and its edges are **marked as used**
3. BFS runs again on the remaining graph to find the next augmenting path
4. This repeats until no more paths exist
5. The resulting set of paths represents the **maximum flow** through the colony

### Optimal Ant Distribution

Not all discovered paths are always beneficial. A longer path may increase the total turns if the ant count does not justify using it. The solver evaluates each subset of paths and calculates the number of turns using:

```
turns = longest_path_length + (N - number_of_paths)
```

It selects the path set that **minimises total turns**, then distributes ants greedily across those paths.

### Summary

| Step | Method |
|------|--------|
| Find paths | BFS (Edmonds-Karp style) |
| Select optimal paths | Turn calculation per subset |
| Distribute ants | Greedy assignment |
| Simulate movement | Turn-by-turn queue |

---

## Project Structure

```
lem-in/
├── main.go            // Entry point
├── graph.go           // Room and tunnel data structures ✅
├── graph_test.go      // Unit tests for graph ✅
├── parser.go          // Input parsing and validation ✅
├── parser_test.go     // Unit tests for parser ✅
├── pathfinder.go      // BFS + augmenting paths
├── solver.go          // Ant distribution and turn optimisation
└── simulator.go       // Turn-by-turn movement output
```

---

## Input Format

The input file must follow this structure:

```
number_of_ants
##start
room_name x y
...
##end
room_name x y
...
room_name x y
...
name1-name2
...
```

- Rooms are defined as `name coord_x coord_y`
- Room names must **not** start with `L` or `#` and must contain **no spaces**
- Tunnels are defined as `name1-name2`
- Lines beginning with `#` (but not `##start` or `##end`) are treated as comments
- `##start` and `##end` are the only valid special commands

---

## Output Format

The program prints the original file content followed by each turn on a new line:

```
number_of_ants
the_rooms
the_links

L1-roomA L2-roomB
L1-roomC L2-roomA L3-roomB
...
```

Where `Lx-y` means ant number `x` moved to room `y`.

---

## Usage

```bash
$ go run . <input_file>
```

Example:

```bash
$ go run . example00.txt
```

---

## Examples

> ⚠️ This section will be populated with verified outputs as examples are confirmed.

### example00 — 4 ants, simple colony

```
$ go run . example00.txt
4
##start
0 0 3
...
L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

Expected: **6 turns or fewer**

---

## Error Handling

The program returns a descriptive error to stdout for all invalid inputs:

```
ERROR: invalid data format
```

With specific messages for each case:

```
ERROR: invalid data format, invalid number of ants
ERROR: invalid data format, no start room found
ERROR: invalid data format, no end room found
ERROR: invalid data format, multiple start rooms found
ERROR: invalid data format, multiple end rooms found
ERROR: invalid data format, duplicate room: <name>
ERROR: invalid data format, unknown room in link: <name>
ERROR: invalid data format, invalid room name: <name>
ERROR: invalid data format, invalid coordinates for room: <name>
ERROR: invalid data format, invalid tunnel format: <tunnel>
ERROR: invalid data format, ##start declared but no room followed
ERROR: invalid data format, ##end declared but no room followed
ERROR: invalid data format, no tunnels found
ERROR: invalid data format, file is empty
ERROR: invalid data format, path is a directory not a file
ERROR: invalid data format, could not open file
ERROR: invalid data format, file read error
```

---

## Testing

```bash
$ go test ./...
```

To see each test name:

```bash
$ go test -v ./...
```

| File | Tests | Status |
|------|-------|--------|
| `graph_test.go` | 8 | ✅ All passing |
| `parser_test.go` | 20 | ✅ All passing |
| `pathfinder_test.go` | — | ⏳ Pending |
| `solver_test.go` | — | ⏳ Pending |
| `simulator_test.go` | — | ⏳ Pending |

---

## Constraints

- Written entirely in **Go**
- Only **standard Go packages** are allowed
- Must handle all edge cases without crashing