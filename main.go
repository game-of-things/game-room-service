package main

import (
	"errors"
	"math/rand"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	Code    string
	Players []Player
}

type Player struct {
	Name string
}

var rooms []Room = make([]Room, 0)

func main() {
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
		room, err := lookupRoom(code)

		if err != nil {
			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, "Room code "+code+" does not exist")
		}
	})

	router.POST("/room/:code/join", func(c *gin.Context) {
		code := c.Param("code")
		room, err := lookupRoom(code)

		if err != nil {
			var player Player

			if err := c.ShouldBindJSON(&player); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			room.Players = append(room.Players, player)

			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, "Room code "+code+" does not exist")
		}
	})

	/*router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})*/

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/rooms", func(c *gin.Context) {
			c.JSON(http.StatusOK, rooms)
		})
	}

	router.Run(":8080")
}

func lookupRoom(code string) (*Room, error) {
	for _, room := range rooms {
		log.Debug("Room code ", room.Code)
		if room.Code == code {
			return &room, nil
		}
	}

	return &Room{}, errors.New("Room does not exist")
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
