package rooms

import "github.com/prometheus/client_golang/prometheus"

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
