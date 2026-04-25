package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/workers"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade to WebSocket connection")
		util.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	workers.RunWebSocketWorker(r.Context(), conn)
}
