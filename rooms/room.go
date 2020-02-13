package rooms

import (
	"errors"
	"math/rand"
	_ "time"

	log "github.com/sirupsen/logrus"
)

var rooms []Room = make([]Room, 0)

// ListRooms list all the available rooms
func ListRooms() *[]Room {
	return &rooms
}

// CreateRoom create a room with a Player resource
func CreateRoom(player Player) *Room {
	room := createRoom()
	room.Players = append(room.Players, player)
	rooms = append(rooms, *room)

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

	return &Room{code, make([]Player, 0)}
}

// LookupRoom find a room by the character code
func LookupRoom(code string) (*Room, error) {
	for index := range rooms {
		log.Debug("Room code ", rooms[index].Code)
		if rooms[index].Code == code {
			log.Debug("Found room: " + code)
			return &rooms[index], nil
		}
	}

	return &Room{}, errors.New("Room code " + code + " does not exist")
}

// AddPlayer add a player to a specified room
func AddPlayer(player Player, room *Room) {
	room.Players = append(room.Players, player)
}
