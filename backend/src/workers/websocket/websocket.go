package websocket

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
	webSocketConnectionCount = atomic.Int64{}
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

func authenticateConnection(ctx context.Context, conn *websocket.Conn) (context.Context, error) {
	authMsg, err := readAuthFrame(conn)
	if err != nil {
		closeUnauthorized(conn, err.Error())
		return ctx, err
	}

	userClaims, err := auth.ValidateJwt(authMsg.Token)
	if err != nil {
		closeUnauthorized(conn, "invalid token")
		return ctx, err
	}

	if err := writeAuthAck(conn); err != nil {
		return ctx, err
	}

	return middleware.SetUserContext(ctx, userClaims), nil
}

func RunWebSocketWorker(ctx context.Context, conn *websocket.Conn) {
	webSocketConnectionCount.Add(1)
	defer func() {
		webSocketConnectionCount.Add(-1)
		_ = conn.Close()
	}()

	logger := log.FromContext(ctx).With().Str("remoteAddr", conn.RemoteAddr().String()).Str("component", "websocket_worker").Logger()
	ctx = log.WithContext(ctx, &logger)

	logger.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection established")

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	authCtx, err := authenticateConnection(ctx, conn)
	if err != nil {
		logger.Error().Err(err).Msg("WebSocket authentication failed")
		return
	}

	logger.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection authenticated successfully")

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
			logger.Warn().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection closed due to context cancellation")
			return
		case <-disconnected:
			logger.Warn().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection closed by client")
			return
		case runMsg := <-monitorRunChannel:
			logger.Trace().Msg("Received monitor run event, sending notification to WebSocket client")
			notification := newMonitorRunNotification(runMsg.Result, runMsg.Monitor)
			notificationBytes, err := json.Marshal(notification)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal monitor run notification")
				continue
			}
			writeMu.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, notificationBytes); err != nil {
				logger.Error().Err(err).Msg("Failed to write monitor run notification to WebSocket connection")
				writeMu.Unlock()
				return
			}
			writeMu.Unlock()
		}
	}
}
