package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/m-milek/leszmonitor/platform/httpx"
	websocketworker "github.com/m-milek/leszmonitor/workers/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	websocketworker.RunWebSocketWorker(ctx, conn)
}
