package main

import (
	"game-room-service/router"
)

func main() {
	router := router.SetupRouter()

	router.Run(":8080")
}
