package rooms

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Room information
type Room struct {
	Code    string
	Players []Player
	Timer   *prometheus.Timer
}

// Player information
type Player struct {
	Name string
}

// RoomMap boom
type RoomMap struct {
	sync.RWMutex
	items map[string]*Room
}

// Add boom
func (roomMap *RoomMap) Add(code string, room *Room) {
	initializeRoomMap(roomMap)
	roomMap.Lock()
	roomMap.items[code] = room
	roomMap.Unlock()
}

// Get boom
func (roomMap *RoomMap) Get(code string) *Room {
	initializeRoomMap(roomMap)

	return roomMap.items[code]
}

// GetAll boom
func (roomMap *RoomMap) GetAll() []*Room {
	initializeRoomMap(roomMap)

	roomMap.RLock()
	slice := make([]*Room, len(roomMap.items))

	index := 0
	for value := range roomMap.items {
		slice[index] = roomMap.Get(value)
		index++
	}
	roomMap.RUnlock()

	return slice
}

// Remove boom
func (roomMap *RoomMap) Remove(code string) {
	roomMap.Lock()
	delete(roomMap.items, code)
	roomMap.Unlock()
}

func initializeRoomMap(roomMap *RoomMap) {
	if roomMap.items == nil {
		roomMap.Lock()
		roomMap.items = make(map[string]*Room)
		roomMap.Unlock()
	}
}
