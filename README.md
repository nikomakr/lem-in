# lem-in

A Go program that solves an ant colony pathfinding problem. Parses a graph of rooms and tunnels, finds optimal non-overlapping paths using DFS and exhaustive vertex-disjoint path selection, and moves N ants from `##start` to `##end` in the fewest turns possible.

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
- [Constraints](#constraints)

---

## Overview

`lem-in` simulates a digital ant farm. Given a colony of rooms connected by tunnels, the program finds the most efficient way to move **N ants** from the `##start` room to the `##end` room in as few turns as possible.

Rules:
- Each room holds at most **one ant at a time** (except `##start` and `##end`)
- Each tunnel can only be used **once per turn**
- All ants begin in `##start` and must reach `##end`

---

## Algorithm

### Path Finding — DFS + Exhaustive Vertex-Disjoint Selection

The program uses a two-stage approach:

**Stage 1 — Find all simple paths (DFS)**

A depth-first search finds every simple path from `##start` to `##end`. A simple path visits each room at most once. Neighbours are sorted alphabetically at each step for deterministic results.

**Stage 2 — Find the best vertex-disjoint set (exhaustive search)**

Two paths are vertex-disjoint if they share no intermediate rooms. The program tries every combination of vertex-disjoint paths and picks the set that minimises total turns for N ants.

This approach guarantees correctness regardless of graph topology, including graphs with cross-edges that would confuse simpler algorithms.

### Why Not Edmonds-Karp / Max Flow?

An Edmonds-Karp BFS-based max flow approach was initially implemented but abandoned because BFS ordering caused incorrect path discovery when cross-edges existed between intended vertex-disjoint paths. The DFS exhaustive approach is correct in all cases.

### Ant Distribution — Greedy Assignment

Once the optimal path set is selected, ants are assigned greedily — each ant goes to whichever path would finish it soonest. This correctly handles unequal path lengths and minimises total turns.

### Turn Simulation

Each ant departs on turn equal to its position in the queue on its path (0-indexed). A room occupancy map enforces the one-ant-per-room rule each turn, with `##end` exempt per spec.

### Summary

| Step | Method |
|------|--------|
| Find all paths | DFS (depth-first search) |
| Select optimal paths | Exhaustive vertex-disjoint combination search |
| Distribute ants | Greedy finish-time assignment |
| Simulate movement | Turn-by-turn with occupancy enforcement |

---

## Project Structure

```
lem-in/
├── main.go              // Entry point ✅
├── graph.go             // Room and tunnel data structures ✅
├── graph_test.go        // Unit tests for graph ✅
├── parser.go            // Input parsing and validation ✅
├── parser_test.go       // Unit tests for parser ✅
├── pathfinder.go        // DFS path finding and disjoint set selection ✅
├── pathfinder_test.go   // Unit tests for pathfinder ✅
├── solver.go            // Ant distribution and turn simulation ✅
├── solver_test.go       // Unit tests for solver ✅
├── simulator.go         // Turn-by-turn output formatting ✅
├── simulator_test.go    // Unit tests for simulator ✅
└── check_errors.sh      // Shell script to verify all error cases ✅
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
- `##start` and `##end` are the only valid special commands — any other `##` command is ignored

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

Or build first for faster execution:

```bash
$ go build -o lem-in .
$ ./lem-in <input_file>
```

---

## Examples

### example00 — 4 ants, 6 turns

```bash
$ go run . example00.txt
```

### example01 — 10 ants, 8 turns

```bash
$ go run . example01.txt
```

### example02 — 20 ants, 11 turns

```bash
$ go run . example02.txt
```

### example06 — 100 ants, under 1.5 minutes

```bash
$ time go run . example06.txt
```

### example07 — 1000 ants, under 2.5 minutes

```bash
$ time go run . example07.txt
```

---

## Error Handling

The program returns a descriptive error for all invalid inputs:

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
ERROR: invalid data format, no path found between start and end
```

To verify all error cases:

```bash
$ ./check_errors.sh
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
| `parser_test.go` | 21 | ✅ All passing |
| `pathfinder_test.go` | 11 | ✅ All passing |
| `solver_test.go` | 9 | ✅ All passing |
| `simulator_test.go` | 8 | ✅ All passing |

**Total: 57 tests, all passing.**

---

## Constraints

- Written entirely in **Go**
- Only **standard Go packages** are allowed
- Must handle all edge cases without crashing