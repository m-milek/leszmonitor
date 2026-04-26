package workers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/auth"
	common "github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
)

type wsAuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type wsAuthResponse struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

var (
	WebSocketConnectionCount atomic.Int64 = atomic.Int64{}
)

func closeUnauthorized(conn *websocket.Conn, reason string) {
	_ = conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.ClosePolicyViolation, reason),
		time.Now().Add(2*time.Second),
	)
}

func readAuthFrame(conn *websocket.Conn) (*wsAuthMessage, error) {
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if messageType != websocket.TextMessage {
		return nil, errors.New("first message must be auth")
	}

	var authMsg wsAuthMessage
	if err := json.Unmarshal(message, &authMsg); err != nil {
		return nil, errors.New("invalid auth message")
	}
	if authMsg.Type != "auth" || strings.TrimSpace(authMsg.Token) == "" {
		return nil, errors.New("invalid auth message")
	}

	return &authMsg, nil
}

func writeAuthAck(conn *websocket.Conn) error {
	authAck, err := json.Marshal(wsAuthResponse{Type: "auth", Status: "ok"})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, authAck)
}

func authenticateConnection(ctx context.Context, conn *websocket.Conn) (context.Context, bool) {
	authMsg, err := readAuthFrame(conn)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to read auth message from WebSocket connection")
		closeUnauthorized(conn, err.Error())
		return ctx, false
	}

	userClaims, err := auth.ValidateJwt(authMsg.Token)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to validate JWT token from WebSocket connection")
		closeUnauthorized(conn, "invalid token")
		return ctx, false
	}

	if err := writeAuthAck(conn); err != nil {
		log.Api.Error().Err(err).Msg("Failed to write auth acknowledgment to WebSocket connection")
		return ctx, false
	}

	log.Api.Info().Any("username", userClaims.Username).Msg("WebSocket connection authenticated successfully")

	return middleware.SetUserContext(ctx, userClaims), true
}

func RunWebSocketWorker(ctx context.Context, conn *websocket.Conn) {
	WebSocketConnectionCount.Add(1)
	defer func() {
		WebSocketConnectionCount.Add(-1)
		_ = conn.Close()
	}()

	log.Api.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection established")

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	authCtx, ok := authenticateConnection(ctx, conn)
	if !ok {
		return
	}

	conn.SetReadDeadline(time.Time{})
	conn.SetWriteDeadline(time.Time{})

	monitorRunChannel := common.MonitorRunChannel.Subscribe()
	defer common.MonitorRunChannel.Unsubscribe(monitorRunChannel)

	var writeMu sync.Mutex

	disconnected := make(chan struct{})
	go func() {
		defer close(disconnected)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if string(msg) == "ping" {
				writeMu.Lock()
				conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				writeMu.Unlock()
			}
		}
	}()

	for {
		select {
		case <-authCtx.Done():
			log.Api.Warn().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection closed due to context cancellation")
			return
		case <-disconnected:
			log.Api.Warn().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection closed by client")
			return
		case runMsg := <-monitorRunChannel:
			log.Uptime.Trace().Msg("Received monitor run event, sending notification to WebSocket client")
			notification := map[string]interface{}{
				"type": "monitor_run",
				"data": runMsg,
			}
			notificationBytes, err := json.Marshal(notification)
			if err != nil {
				log.Api.Error().Err(err).Msg("Failed to marshal monitor run notification")
				continue
			}
			writeMu.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, notificationBytes); err != nil {
				log.Api.Error().Err(err).Msg("Failed to write monitor run notification to WebSocket connection")
				writeMu.Unlock()
				return
			}
			writeMu.Unlock()
		}
	}
}
