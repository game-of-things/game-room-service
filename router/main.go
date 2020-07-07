package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"game-room-service/rooms"

	cors "github.com/itsjamie/gin-cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter create the router
func SetupRouter() *gin.Engine {
	log.SetLevel(log.DebugLevel)

	log.Info("Starting game room service")

	router := gin.Default()

	corsMiddleware := cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	})

	router.Use(corsMiddleware)

	router.POST("/rooms/create", func(c *gin.Context) {
		var player rooms.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room := rooms.CreateRoom(&player)

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
			player.Join(room)

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
			if err := player.Quit(room); err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Player " + player.Name + " quit room " + code})
		} else {
			c.JSON(http.StatusNotFound, err.Error())
		}
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router.GET("/room/:code/player/:playerName/websocket", func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Error("Failed to set websocket upgrade: %+v", err)
			return
		}

		for {
			t, _, err := conn.ReadMessage()
			if err != nil {
				log.Error("Error reading message: %+v", err)
				continue
			}
			if t == websocket.BinaryMessage {
				log.Warn("Received binary message when expecting plain text. Ignoring")
				continue
			}
			if t == websocket.PingMessage {
				conn.WriteControl(websocket.PongMessage, []byte("pong"), time.Now())
			}

		}
	})

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/rooms", func(c *gin.Context) {
			c.JSON(http.StatusOK, rooms.ListRooms())
		})
	}

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
