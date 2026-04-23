package main

import "testing"

// --- AddRoom tests ---

func TestAddRoom_Success(t *testing.T) {
	farm := NewFarm()
	room := &Room{Name: "start", X: 1, Y: 2, IsStart: true}

	if !farm.AddRoom(room) {
		t.Error("expected AddRoom to return true for a new room")
	}
	if _, exists := farm.Rooms["start"]; !exists {
		t.Error("expected room 'start' to exist in farm after adding")
	}
}

func TestAddRoom_Duplicate(t *testing.T) {
	farm := NewFarm()
	room := &Room{Name: "start", X: 1, Y: 2}

	farm.AddRoom(room)
	if farm.AddRoom(room) {
		t.Error("expected AddRoom to return false for a duplicate room")
	}
}

func TestAddRoom_MultipleRooms(t *testing.T) {
	farm := NewFarm()
	rooms := []*Room{
		{Name: "a", X: 0, Y: 0},
		{Name: "b", X: 1, Y: 1},
		{Name: "c", X: 2, Y: 2},
	}
	for _, r := range rooms {
		if !farm.AddRoom(r) {
			t.Errorf("expected AddRoom to succeed for room %s", r.Name)
		}
	}
	if len(farm.Rooms) != 3 {
		t.Errorf("expected 3 rooms, got %d", len(farm.Rooms))
	}
}

// --- AddTunnel tests ---

func TestAddTunnel_Success(t *testing.T) {
	farm := NewFarm()
	farm.AddRoom(&Room{Name: "a"})
	farm.AddRoom(&Room{Name: "b"})

	if !farm.AddTunnel("a", "b") {
		t.Error("expected AddTunnel to return true for valid rooms")
	}

	// Check bidirectional
	if len(farm.Tunnels["a"]) != 1 || farm.Tunnels["a"][0] != "b" {
		t.Error("expected tunnel from a to b")
	}
	if len(farm.Tunnels["b"]) != 1 || farm.Tunnels["b"][0] != "a" {
		t.Error("expected tunnel from b to a (bidirectional)")
	}
}

func TestAddTunnel_UnknownRoom(t *testing.T) {
	farm := NewFarm()
	farm.AddRoom(&Room{Name: "a"})

	if farm.AddTunnel("a", "ghost") {
		t.Error("expected AddTunnel to return false when a room does not exist")
	}
}

func TestAddTunnel_Duplicate(t *testing.T) {
	farm := NewFarm()
	farm.AddRoom(&Room{Name: "a"})
	farm.AddRoom(&Room{Name: "b"})
	farm.AddTunnel("a", "b")

	if farm.AddTunnel("a", "b") {
		t.Error("expected AddTunnel to return false for a duplicate tunnel")
	}
}

func TestAddTunnel_SelfLink(t *testing.T) {
	farm := NewFarm()
	farm.AddRoom(&Room{Name: "a"})

	// A room linking to itself should be rejected (duplicate check catches it)
	if farm.AddTunnel("a", "a") {
		t.Error("expected AddTunnel to return false for a self-link")
	}
}

// --- NewFarm tests ---

func TestNewFarm_IsEmpty(t *testing.T) {
	farm := NewFarm()

	if farm.Ants != 0 {
		t.Errorf("expected 0 ants, got %d", farm.Ants)
	}
	if len(farm.Rooms) != 0 {
		t.Errorf("expected empty rooms map, got %d entries", len(farm.Rooms))
	}
	if len(farm.Tunnels) != 0 {
		t.Errorf("expected empty tunnels map, got %d entries", len(farm.Tunnels))
	}
	if farm.StartRoom != "" || farm.EndRoom != "" {
		t.Error("expected StartRoom and EndRoom to be empty strings")
	}
}