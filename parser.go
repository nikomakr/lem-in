package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseFarm reads the input file and returns a populated Farm or an error.
func ParseFarm(filename string) (*Farm, error) {
	// --- File-level validation ---
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, could not open file")
	}
	if info.IsDir() {
		return nil, fmt.Errorf("ERROR: invalid data format, path is a directory not a file")
	}
	if info.Size() == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, file is empty")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, could not open file: %v", err)
	}
	defer file.Close()

	farm := NewFarm()

	// Increase scanner buffer to handle long lines (default is 64KB)
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 1024*1024) // 1MB buffer
	scanner.Buffer(buf, 1024*1024)

	antsParsed := false
	nextIsStart := false
	nextIsEnd := false
	parsingLinks := false
	hasLines := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}
		hasLines = true

		// --- Ant count (must be the very first non-empty line) ---
		if !antsParsed {
			if strings.Contains(line, ".") {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
			}
			ants, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
			}
			if ants <= 0 {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
			}
			farm.Ants = ants
			antsParsed = true
			continue
		}

		// --- Special commands ---
		if line == "##start" {
			if farm.StartRoom != "" {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple start rooms found")
			}
			nextIsStart = true
			continue
		}
		if line == "##end" {
			if farm.EndRoom != "" {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple end rooms found")
			}
			nextIsEnd = true
			continue
		}

		// --- Comments ---
		if strings.HasPrefix(line, "#") {
			continue
		}

		// --- Tunnel definition ---
		if !parsingLinks && strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parsingLinks = true
		}

		if parsingLinks {
			if !strings.Contains(line, "-") {
				continue
			}
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid tunnel format: %s", line)
			}
			if parts[0] == "" || parts[1] == "" {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid tunnel format: %s", line)
			}
			if _, ok := farm.Rooms[parts[0]]; !ok {
				return nil, fmt.Errorf("ERROR: invalid data format, unknown room in link: %s", parts[0])
			}
			if _, ok := farm.Rooms[parts[1]]; !ok {
				return nil, fmt.Errorf("ERROR: invalid data format, unknown room in link: %s", parts[1])
			}
			farm.AddTunnel(parts[0], parts[1])
			continue
		}

		// --- Room definition ---
		parts := strings.Fields(line)

		if len(parts) != 3 {
			if nextIsStart {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid room definition after ##start: %s", line)
			}
			if nextIsEnd {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid room definition after ##end: %s", line)
			}
			continue
		}

		name := parts[0]

		if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
			return nil, fmt.Errorf("ERROR: invalid data format, invalid room name: %s", name)
		}

		if strings.Contains(parts[1], ".") || strings.Contains(parts[2], ".") {
			return nil, fmt.Errorf("ERROR: invalid data format, invalid coordinates for room: %s", name)
		}

		x, errX := strconv.Atoi(parts[1])
		y, errY := strconv.Atoi(parts[2])
		if errX != nil || errY != nil {
			return nil, fmt.Errorf("ERROR: invalid data format, invalid coordinates for room: %s", name)
		}

		room := &Room{
			Name:    name,
			X:       x,
			Y:       y,
			IsStart: nextIsStart,
			IsEnd:   nextIsEnd,
		}

		if !farm.AddRoom(room) {
			return nil, fmt.Errorf("ERROR: invalid data format, duplicate room: %s", name)
		}

		if nextIsStart {
			farm.StartRoom = name
			nextIsStart = false
		}
		if nextIsEnd {
			farm.EndRoom = name
			nextIsEnd = false
		}
	}

	// --- IO / process interruption error ---
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, file read error: %v", err)
	}

	// --- Whitespace-only file ---
	if !hasLines {
		return nil, fmt.Errorf("ERROR: invalid data format, file contains no valid data")
	}

	// --- ##start or ##end with no room following ---
	if nextIsStart {
		return nil, fmt.Errorf("ERROR: invalid data format, ##start declared but no room followed")
	}
	if nextIsEnd {
		return nil, fmt.Errorf("ERROR: invalid data format, ##end declared but no room followed")
	}

	// --- Missing start or end ---
	if farm.StartRoom == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if farm.EndRoom == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, no end room found")
	}

	// --- No tunnels at all ---
	if len(farm.Tunnels) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, no tunnels found")
	}

	return farm, nil
}