package main

import (
	"game-room-service/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/json"
)

func TestPingRoute(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms", nil)
	router.ServeHTTP(w, req)

	json.NewDecoder(w.Body)

	assert.Equal(t, 200, w.Code)

	//assert.Equal(t, "pong", w.Body.String())
}
