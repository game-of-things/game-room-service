package rooms

// Room information
type Room struct {
	Code    string
	Players []Player
}

// Player information
type Player struct {
	Name string
}
