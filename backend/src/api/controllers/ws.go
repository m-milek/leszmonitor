package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
	util "github.com/m-milek/leszmonitor/api/api_util"
	websocket2 "github.com/m-milek/leszmonitor/workers/websocket"
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
		util.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	websocket2.RunWebSocketWorker(ctx, conn)
}
