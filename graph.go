package main

// Room represents a node in the colony graph.
type Room struct {
	Name    string
	X, Y    int
	IsStart bool
	IsEnd   bool
}

// Farm holds the entire parsed colony.
type Farm struct {
	Ants      int
	Rooms     map[string]*Room
	Tunnels   map[string][]string // adjacency list: room name -> list of connected room names
	StartRoom string
	EndRoom   string
}

// NewFarm initialises an empty Farm.
func NewFarm() *Farm {
	return &Farm{
		Rooms:   make(map[string]*Room),
		Tunnels: make(map[string][]string),
	}
}

// AddRoom adds a room to the farm. Returns false if a room with that name already exists.
func (f *Farm) AddRoom(r *Room) bool {
	if _, exists := f.Rooms[r.Name]; exists {
		return false
	}
	f.Rooms[r.Name] = r
	return true
}

// AddTunnel adds a bidirectional tunnel between two rooms.
// Returns false if either room does not exist, the tunnel already exists,
// or the tunnel links a room to itself.
func (f *Farm) AddTunnel(a, b string) bool {
	if a == b {
		return false
	}
	if _, ok := f.Rooms[a]; !ok {
		return false
	}
	if _, ok := f.Rooms[b]; !ok {
		return false
	}
	// Check for duplicate tunnel
	for _, neighbour := range f.Tunnels[a] {
		if neighbour == b {
			return false
		}
	}
	f.Tunnels[a] = append(f.Tunnels[a], b)
	f.Tunnels[b] = append(f.Tunnels[b], a)
	return true
}