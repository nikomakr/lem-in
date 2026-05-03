package main

import (
	"fmt"
	"testing"
)

func TestDebugExample01Paths(t *testing.T) {
	farm := NewFarm()
	farm.Ants = 10
	farm.StartRoom = "start"
	farm.EndRoom = "end"

	rooms := []string{"start", "end", "0", "o", "n", "e", "t", "E", "a", "m", "h", "A", "c", "k"}
	for _, name := range rooms {
		r := &Room{Name: name}
		if name == "start" {
			r.IsStart = true
		}
		if name == "end" {
			r.IsEnd = true
		}
		farm.AddRoom(r)
	}

	tunnels := [][2]string{
		{"start", "t"}, {"start", "h"}, {"start", "0"},
		{"t", "E"}, {"E", "a"}, {"a", "m"}, {"m", "end"},
		{"h", "A"}, {"A", "c"}, {"c", "k"}, {"k", "end"},
		{"0", "o"}, {"o", "n"}, {"n", "e"}, {"e", "end"},
		{"n", "m"}, {"h", "n"}, {"n", "o"}, {"m", "a"},
	}
	for _, tun := range tunnels {
		farm.AddTunnel(tun[0], tun[1])
	}

	capacity := buildCapacityGraph(farm)

	// First BFS
	path1 := bfsCapacity(farm.StartRoom+"_in", farm.EndRoom+"_out", capacity)
	fmt.Printf("Path 1: %v\n", toRealPath(path1))

	// Update capacities
	for i := 0; i < len(path1)-1; i++ {
		capacity[path1[i]][path1[i+1]]--
		capacity[path1[i+1]][path1[i]]++
	}

	// Second BFS
	path2 := bfsCapacity(farm.StartRoom+"_in", farm.EndRoom+"_out", capacity)
	if path2 == nil {
		fmt.Println("Path 2: nil — checking why...")
		// Check what start_in can reach
		fmt.Println("start_in capacity:")
		for k, v := range capacity["start_in"] {
			fmt.Printf("  start_in -> %s = %d\n", k, v)
		}
		fmt.Println("start_out capacity:")
		for k, v := range capacity["start_out"] {
			fmt.Printf("  start_out -> %s = %d\n", k, v)
		}
	} else {
		fmt.Printf("Path 2: %v\n", toRealPath(path2))
	}
}