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
	/*if os.Getenv("PRODUCTION") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else if os.Getenv("") == "" {

	} else {*/
	log.SetLevel(log.DebugLevel)
	gin.SetMode(gin.TestMode)

	log.Info("Starting game room service")

	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	router.POST("/room/create", func(c *gin.Context) {
		room := createRoom()

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

	router.GET("/room/:code/join", func(c *gin.Context) {
		
	})

	if gin.Mode() == gin.TestMode {
		router.GET("/room", func(c *gin.Context) {
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

	for i := 0; i < len(runes); i++ {
		runes[i] = rune(rand.Intn(26) + 65)
	}

	return &Room{string(runes), make([]Player, 0)}
}
