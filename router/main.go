package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"game-room-service/rooms"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter create the router
func SetupRouter() *gin.Engine {
	log.SetLevel(log.DebugLevel)

	log.Info("Starting game room service")

	router := gin.Default()

	router.POST("/rooms/create", func(c *gin.Context) {
		var player rooms.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room := rooms.CreateRoom(player)

		c.JSON(http.StatusOK, room)
	})

	router.GET("/room/:code", func(c *gin.Context) {
		code := c.Param("code")

		if room, err := rooms.LookupRoom(code); err == nil {
			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.POST("/room/:code/join", func(c *gin.Context) {
		var player rooms.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code := c.Param("code")

		log.Debug("Attempting to join room" + code)

		if room, err := rooms.LookupRoom(code); err == nil {
			rooms.Join(player, room)

			c.JSON(http.StatusOK, room)
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.POST("/room/:code/quit", func(c *gin.Context) {
		var player rooms.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code := c.Param("code")

		log.Debug("Attempting to quit room " + code)

		if room, err := rooms.LookupRoom(code); err == nil {
			if err := rooms.Quit(player, room); err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Player quit"})
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/rooms", func(c *gin.Context) {
			c.JSON(http.StatusOK, rooms.ListRooms())
		})
	}

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
