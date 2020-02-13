package main

import (
	"game-room-service/router"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	router := router.SetupRouter()

	router.Run(":8080")

}
