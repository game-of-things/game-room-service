package rooms

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	//rooms []Room = make([]Room, 0)

	roomsMap = RoomMap{}

	activeRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gameofthings_active_rooms_total",
		Help: "Total number of active rooms",
	})

	roomDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "gameofthings_room_duration_seconds",
		Help: "Total duration of active rooms",
	})
)

func init() {
	prometheus.MustRegister(activeRooms)
	prometheus.MustRegister(roomDuration)
}

// ListRooms list all the available rooms
func ListRooms() []*Room {
	return roomsMap.GetAll()
}

// CreateRoom create a room with a Player resource
func CreateRoom(player *Player) *Room {
	room := createRoom()
	room.Players = append(room.Players, player)
	room.Timer = prometheus.NewTimer(roomDuration)

	roomsMap.Add(room.Code, room)

	activeRooms.Inc()

	log.Debug(&roomsMap)

	return room
}

func createRoom() *Room {
	runes := make([]rune, 4)

	var code string

	for {
		for i := 0; i < len(runes); i++ {
			runes[i] = rune(rand.Intn(26) + 65)
		}

		code = string(runes)
		_, err := LookupRoom(code)

		if err != nil {
			break
		}
	}

	return &Room{code, make([]*Player, 0), nil}
}

// LookupRoom find a room by the character code
func LookupRoom(code string) (*Room, error) {
	if room := roomsMap.Get(code); room != nil {
		return room, nil
	}

	return &Room{}, errors.New("Room code " + code + " does not exist")
}

// Join add a player to a specified room
func (player *Player) Join(room *Room) {
	for _, other := range room.Players {
		other.Connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"newPlayer\":\"%v\"", player.Name)))
	}
	room.Players = append(room.Players, player)
}

// LookupPlayer lookup player by name
func (room *Room) LookupPlayer(playerName string) (*Player, error) {
	for index := range room.Players {
		if room.Players[index].Name == playerName {
			return room.Players[index], nil
		}
	}

	return &Player{}, errors.New("Player " + playerName + " not found in room " + room.Code)
}

// Quit remove player from specified room
func (player *Player) Quit(room *Room) error {
	for index := range room.Players {
		if room.Players[index].Name == player.Name {
			if len(room.Players) <= 1 {
				room.remove()
				return nil
			}
			room.Players = append(room.Players[:index], room.Players[index+1:]...)
			return nil
		}
	}

	return errors.New("Player " + player.Name + " not found in room " + room.Code)
}

func (room *Room) remove() {
	log.Debug("Removing room " + room.Code)
	roomsMap.Remove(room.Code)
	room.Timer.ObserveDuration()
	activeRooms.Dec()
}
