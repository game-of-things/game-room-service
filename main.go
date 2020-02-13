package main

import (
	"game-room-service/router"
	"math/rand"
	"time"
)

func main() {

	// set random seed for generating unique values
	rand.Seed(time.Now().UnixNano())

	router := router.SetupRouter()

	router.Run(":8080")

}
