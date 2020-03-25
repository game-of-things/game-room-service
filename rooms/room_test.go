package rooms

import (
	"regexp"
	"testing"
)

var (
	validRoomCode = regexp.MustCompile(`^[A-Z]{4}$`)
)

func TestCreateRoomCode(t *testing.T) {
	player := Player{}

	room := CreateRoom(player)

	if validRoomCode.Match([]byte(room.Code)) {
		t.Errorf("Room code was invalid: %v", room.Code)
	}
}
