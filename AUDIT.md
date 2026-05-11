# lem-in — Audit Checklist

This document maps every audit question to the expected answer and
how to verify it. Work through each item in order before your audit session.

---

## Functional

---

### Has the requirement for allowed packages been respected?

**Expected:** Only standard Go packages used.

**How to verify:**
```bash
grep -A 10 "import" parser.go
grep -A 5 "import" main.go
grep -A 5 "import" simulator.go
grep -A 3 "import" pathfinder.go
```

**What you should see:** Only packages from the Go standard library:
- `bufio`, `fmt`, `os`, `strconv`, `strings` in `parser.go`
- `sort` in `pathfinder.go`
- `fmt`, `strings` in `simulator.go`
- `fmt`, `os` in `main.go`

No third-party imports anywhere.

---

### Is the program able to read the ant farm from a file argument?

**How to verify:**
```bash
go run . example00.txt
```

**Expected:** Program reads and prints the file content followed by moves.

---

### Does the program accept only ##start and ##end as special commands?

**How to verify:**
Try to use ##start twice — should error:
```bash
echo "2
##start
a 0 0
##start
b 1 1
##end
c 2 2
a-b
b-c" > /tmp/test_double_start.txt && go run . /tmp/test_double_start.txt
```
expected outcome:
```
ERROR: invalid data format, multiple start rooms found
```

# Try to use ##end twice — should error
```bash
echo "2
##start
a 0 0
##end
b 1 1
##end
c 2 2
a-b
b-c" > /tmp/test_double_end.txt && go run . /tmp/test_double_end.txt
```

expected outcome:
```
ERROR: invalid data format, multiple end rooms found
```

# Any other ##something is not a special command — treated as a comment
```bash
echo "2
##start
a 0 0
##end
b 1 1
##whatever
a-b" > /tmp/test_unknown.txt && go run . /tmp/test_unknown.txt
```
expected outcome:
```
2
##start
a 0 0
##end
b 1 1
##whatever
a-b

L1-b
L2-b
```

**Check in `parser.go`:**
```go
if line == "##start" { ... }
if line == "##end" { ... }
```

Anything else starting with `#` is treated as a comment and skipped. As long as an input/test has ##start & ##end then the programme runs.

---

### Are moves printed in the correct format?

**Format required:**
```
Lx-y Lx-y Lx-y
```
One line per turn. Each move is `Lx-y` where `x` is the ant number and `y` is the room name.

**How to verify:**
```bash
go run . example00.txt
```

**Implemented in:** `simulator.go` → `formatTurn()` 

---

### example00 — at most 6 turns

**Run:**
```bash
go run . example00.txt
```

**Expected output:**
```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

**Verify:** Count the move lines — must be 6 or fewer. 

---

### example01 — at most 8 turns

**Run:**
```bash
go run . example01.txt
```

**Expected output moves (8 turns):**
```
L1-t L2-h L3-0
L1-E L2-A L3-o L4-t L5-h L6-0
L1-a L2-c L3-n L4-E L5-A L6-o L7-t L8-h L9-0
L1-m L2-k L3-e L4-a L5-c L6-n L7-E L8-A L9-o L10-t
L1-end L2-end L3-end L4-m L5-k L6-e L7-a L8-c L9-n L10-E
L4-end L5-end L6-end L7-m L8-k L9-e L10-a
L7-end L8-end L9-end L10-m
L10-end
```

**Verify:** Count move lines — must be 8 or fewer. 

---

### example02 — at most 11 turns

**Run:**
```bash
go run . example02.txt
```

**Verify:** Count move lines — must be 11 or fewer.

**Run:**
```bash
go build -o lem-in . && ./lem-in example02.txt | grep "^L" | wc -l
```

---

### example03 — at most 6 turns

**Run:**
```bash
go run . example03.txt
```

**Verify:** Count move lines — must be 6 or fewer. 

---

### example04 — at most 6 turns

**Run:**
```bash
go run . example04.txt
```

**Verify:** Count move lines — must be 6 or fewer. 

---

### example05 — at most 8 turns

**Run:**
```bash
go run . example05.txt
```

**Verify:** Count move lines — must be 8 or fewer.

---

### badexample00 — ERROR output

**Run:**
```bash
go run . badexample00.txt
```

**Expected:**
```
ERROR: invalid data format
```
or more specific. Our output:
```
ERROR: invalid data format, invalid number of ants
```

---

### badexample01 — ERROR output

**Run:**
```bash
go run . badexample01.txt
```

**Expected:**
```
ERROR: invalid data format
```
or more specific. Our output:
```
ERROR: invalid data format, no path found between start and end
```

---

### example06 — 100 ants, under 1.5 minutes

**Run:**
```bash
time go run . example06.txt
```

**Result (example):** ~go run . example06.txt  0.05s user 0.11s system 63% cpu 0.253 total

---

### example07 — 1000 ants, under 2.5 minutes

**Run:**
```bash
time go run . example07.txt
```

**Result:** ~go run . example07.txt  0.06s user 0.10s system 60% cpu 0.266 total

---

### Are ants alone in each room per turn?

**How to verify manually:** For each turn line, check no room name appears twice.

**Example check script:**
```bash
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt example06.txt example07.txt; do
    echo -n "$f: "
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    startroom=$(awk '/##start/{found=1; next} found{print $1; exit}' $f)
    go run . $f | grep "^L" | awk -v endroom="$endroom" -v startroom="$startroom" '{
        delete seen
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            room = parts[2]
            if (room == endroom || room == startroom) continue
            if (room in seen) print "DUPLICATE ROOM: " room
            seen[room] = 1
        }
    } END { print "OK" }'
done
```

Each intermediate room holds at most one ant per turn. Start and end are exempt.

---

### Is each tunnel used only once per turn?

**How to verify:** 
**Example check script:**
```bash
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt example06.txt example07.txt; do
    echo -n "$f: "
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    startroom=$(awk '/##start/{found=1; next} found{print $1; exit}' $f)
    go run . $f | grep "^L" | awk -v endroom="$endroom" -v startroom="$startroom" '
    NR > 0 {
        delete seen
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            ant = substr(parts[1], 2)
            room = parts[2]
            tunnel = (prevRoom[ant] < room) ? prevRoom[ant] "-" room : room "-" prevRoom[ant]
            if (tunnel in seen && prevRoom[ant] != "" && prevRoom[ant] != startroom) {
                print "DUPLICATE TUNNEL: " tunnel " in turn: " $0
            }
            seen[tunnel] = 1
            prevRoom[ant] = room
        }
    }
    END { print "OK" }'
done
```

Each ant moves one step per turn along its assigned path.
Paths are non-overlapping (guaranteed by Edmonds-Karp node splitting).
No two ants share a tunnel in the same turn.

---

### Are all ants in ##end at the finish?

**How to verify:** The last move line for each ant must end at the end room.

```bash
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt example06.txt example07.txt; do
    echo -n "$f: "
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    numants=$(head -1 $f)
    go run . $f | grep "^L" | awk -v endroom="$endroom" -v numants="$numants" '
    {
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            ant = parts[1]
            room = parts[2]
            lastroom[ant] = room
        }
    }
    END {
        failed = 0
        for (ant in lastroom) {
            if (lastroom[ant] != endroom) {
                print "FAIL: " ant " last seen in " lastroom[ant] " not " endroom
                failed = 1
            }
        }
        if (length(lastroom) != numants) {
            print "FAIL: only " length(lastroom) " of " numants " ants moved"
            failed = 1
        }
        if (!failed) print "OK - all " numants " ants reached " endroom
    }'
done
```

This proves three things at once:

Every ant's last recorded room is the end room 
Every single ant moved at least once 
The total number of ants that moved matches the declared ant count

---

### Are results always correct across multiple runs?

**How to verify:**
```bash
for i in 1 2 3 4 5; do go run . example01.txt | grep "^L" | wc -l; done
```

All runs produce the same deterministic output.

---

### Does the program handle all error cases with a message?

```bash
go build -o lem-in . && ./check_errors.sh
```

**Full list of error messages produced:**
```
✅ ERROR: invalid data format, invalid number of ants
✅ ERROR: invalid data format, no start room found
✅ ERROR: invalid data format, no end room found
✅ ERROR: invalid data format, multiple start rooms found
✅ ERROR: invalid data format, multiple end rooms found
✅ ERROR: invalid data format, duplicate room: a
✅ ERROR: invalid data format, unknown room in link: ghost
✅ ERROR: invalid data format, invalid room name: L1
✅ ERROR: invalid data format, invalid coordinates for room: a
✅ ERROR: invalid data format, invalid tunnel format: a-b-c
✅ ERROR: invalid data format, ##start declared but no room followed
✅ ERROR: invalid data format, ##end declared but no room followed
✅ ERROR: invalid data format, no tunnels found
✅ ERROR: invalid data format, file is empty
✅ ERROR: invalid data format, path is a directory not a file
✅ ERROR: invalid data format, could not open file
✅ ERROR: invalid data format, file contains no valid data
✅ ERROR: invalid data format, no path found between start and end
```

---

### Does the solution move ants from ##start to ##end properly?

```bash
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt example06.txt example07.txt; do
    echo -n "$f: "
    startroom=$(awk '/##start/{found=1; next} found{print $1; exit}' $f)
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    numants=$(head -1 $f)
    go run . $f | grep "^L" | awk -v start="$startroom" -v end="$endroom" -v numants="$numants" '
    {
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            ant = parts[1]
            room = parts[2]
            if (!(ant in firstroom)) firstroom[ant] = room
            lastroom[ant] = room
        }
    }
    END {
        failed = 0
        if (length(lastroom) != numants) {
            print "FAIL: only " length(lastroom) " of " numants " ants moved"
            failed = 1
        }
        for (ant in lastroom) {
            if (lastroom[ant] != end) {
                print "FAIL: " ant " ended in " lastroom[ant] " not " end
                failed = 1
            }
        }
        if (!failed) print "OK — " numants " ants moved from " start " to " end
    }'
done
```

---

## General (Bonus)

---

### ➕ Does the program present an ant farm visualizer?

Not implemented. This is a bonus feature beyond the project requirements.

---

### ➕ Is it possible to see the ants moving?

Not implemented. Follows from the visualizer not being present.

---

### ➕ Is the error output more specific?

Yes 
— run : 
./check_errors.sh 
to prove it. Every error includes a specific reason beyond the base message, for example ERROR: invalid data format, duplicate room: a rather than just ERROR: invalid data format.

---

### ➕ Is the visualizer in 3D?

Not implemented.

---

## Basic (Bonus)

---

### ➕ Does the project run quickly and effectively?

Yes  — run:
```bash
go build -o lem-in .
time ./lem-in example06.txt > /dev/null
time ./lem-in example07.txt > /dev/null
```

100 ants completes in under 0.25 seconds. 1000 ants completes in under 0.03 seconds. Both are hundreds of times faster than the audit limits.

---

### ➕ Is there a test file for this code?

Yes — run:
```bash
go test -v ./...
```
57 tests across 5 test files covering every module.

---

### ➕ Are the tests checking each possible case?

Yes — tests cover:

- Valid input happy paths
- All 18 error cases via check_errors.sh
- Edge cases: empty files, directories, self-links, duplicate rooms, duplicate tunnels, float coordinates, zero ants, missing start/end
- Algorithm correctness: single path, parallel paths, conflicting paths, three disjoint paths
- Simulator: single ant, multiple ants, empty assignment
Output format: single move, multiple moves, named rooms, audit format match

---

### ➕ Does the code obey good practices?


Yes:

- Every exported function has a comment explaining its purpose
- Clear separation of concerns — parsing, path finding, solving and output are in separate files
- No global state
- Errors are handled explicitly at every step
- No unnecessary data requests or repeated computation
- Deterministic output via sorted neighbour traversal

---

## How to Run All Verifications at Once

# Build
go build -o lem-in .

# All unit tests
go test -v ./...

# Turn counts — must be within limits
./lem-in example00.txt | grep "^L" | wc -l   # ≤ 6
./lem-in example01.txt | grep "^L" | wc -l   # ≤ 8
./lem-in example02.txt | grep "^L" | wc -l   # ≤ 11
./lem-in example03.txt | grep "^L" | wc -l   # ≤ 6
./lem-in example04.txt | grep "^L" | wc -l   # ≤ 6
./lem-in example05.txt | grep "^L" | wc -l   # ≤ 8

# Error cases — all 18 must show ✅
./check_errors.sh

# Determinism — all runs must show same number
for i in 1 2 3 4 5; do ./lem-in example01.txt | grep "^L" | wc -l; done

# No duplicate rooms per turn
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt; do
    echo -n "$f: "
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    startroom=$(awk '/##start/{found=1; next} found{print $1; exit}' $f)
    ./lem-in $f | grep "^L" | awk -v endroom="$endroom" -v startroom="$startroom" '{
        delete seen
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            room = parts[2]
            if (room == endroom || room == startroom) continue
            if (room in seen) print "DUPLICATE ROOM: " room
            seen[room] = 1
        }
    } END { print "OK" }'
done

# All ants reach end
for f in example00.txt example01.txt example02.txt example03.txt example04.txt example05.txt example06.txt example07.txt; do
    echo -n "$f: "
    endroom=$(awk '/##end/{found=1; next} found{print $1; exit}' $f)
    numants=$(head -1 $f)
    ./lem-in $f | grep "^L" | awk -v end="$endroom" -v numants="$numants" '
    {
        n = split($0, moves, " ")
        for (i = 1; i <= n; i++) {
            split(moves[i], parts, "-")
            lastroom[parts[1]] = parts[2]
        }
    }
    END {
        failed = 0
        if (length(lastroom) != numants) { print "FAIL: only " length(lastroom) " of " numants " ants moved"; failed = 1 }
        for (ant in lastroom) if (lastroom[ant] != end) { print "FAIL: " ant " ended in " lastroom[ant]; failed = 1 }
        if (!failed) print "OK — " numants " ants reached " end
    }'
done

# Performance
time ./lem-in example06.txt > /dev/null
time ./lem-in example07.txt > /dev/null
