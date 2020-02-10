package main

import (
	"errors"
	"math/rand"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var rooms []Room = make([]Room, 0)

func main() {
	router := setupRouter()

	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	log.SetLevel(log.DebugLevel)

	log.Info("Starting game room service")

	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	router.POST("/rooms/create", func(c *gin.Context) {
		var player Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room := createRoom()
		room.Players = append(room.Players, player)
		rooms = append(rooms, *room)

		c.JSON(http.StatusOK, room)
	})

	router.GET("/room/:code", func(c *gin.Context) {
		code := c.Param("code")

		if room, err := lookupRoom(code); err == nil {
			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.POST("/room/:code/join", func(c *gin.Context) {
		code := c.Param("code")

		log.Debug("Attempting to join room" + code)

		if room, err := lookupRoom(code); err == nil {
			var player Player

			if err := c.ShouldBindJSON(&player); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			room.Players = append(room.Players, player)

			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/rooms", func(c *gin.Context) {
			c.JSON(http.StatusOK, rooms)
		})
	}

	return router
}

func lookupRoom(code string) (*Room, error) {
	for index := range rooms {
		log.Debug("Room code ", rooms[index].Code)
		if rooms[index].Code == code {
			log.Debug("Found room: " + code)
			return &rooms[index], nil
		}
	}

	return &Room{}, errors.New("Room code " + code + " does not exist")
}

func createRoom() *Room {
	runes := make([]rune, 4)

	var code string

	for {
		for i := 0; i < len(runes); i++ {
			runes[i] = rune(rand.Intn(26) + 65)
		}

		code = string(runes)
		_, err := lookupRoom(code)

		if err != nil {
			break
		}
	}

	return &Room{code, make([]Player, 0)}
}
