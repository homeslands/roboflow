package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotService "github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotCommandService "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type RaybotClient struct {
	raybot *model.Raybot

	cmdMap  map[string]*model.RaybotCommand
	cmdChan chan *model.RaybotCommand

	// services
	raybotSvc raybotService.Service
	cmdSvc    raybotCommandService.Service

	conn   *websocket.Conn
	send   chan []byte
	logger *slog.Logger
}

func NewRaybotClient(
	raybot *model.Raybot,
	raybotSvc raybotService.Service,
	cmdSvc raybotCommandService.Service,
	conn *websocket.Conn,
	logger *slog.Logger,
) *RaybotClient {
	return &RaybotClient{
		raybot:    raybot,
		cmdMap:    make(map[string]*model.RaybotCommand),
		cmdChan:   make(chan *model.RaybotCommand, 1),
		raybotSvc: raybotSvc,
		cmdSvc:    cmdSvc,
		conn:      conn,
		send:      make(chan []byte, 256),
		logger:    logger,
	}
}

func InitilizeRaybotClient(
	ctx context.Context,
	raybotID uuid.UUID,
	raybotSvc raybotService.Service,
	cmdSvc raybotCommandService.Service,
	conn *websocket.Conn,
	logger *slog.Logger,
) (*RaybotClient, error) {
	raybot, err := raybotSvc.GetByID(ctx, raybotService.GetRaybotByIDQuery{ID: raybotID})
	if err != nil {
		return nil, fmt.Errorf("error getting raybot: %w", err)
	}

	err = raybotSvc.UpdateState(ctx, raybotService.UpdateStateCommand{
		ID:    raybotID,
		State: model.RaybotStateIdle,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating raybot status: %w", err)
	}

	return NewRaybotClient(
		raybot,
		raybotSvc,
		cmdSvc,
		conn,
		logger.With(slog.String("raybot_id", raybotID.String())),
	), nil
}

func (c RaybotClient) ID() string {
	return c.raybot.ID.String()
}

func (c *RaybotClient) Close() {
	err := c.raybotSvc.UpdateState(context.Background(), raybotService.UpdateStateCommand{
		ID:    c.raybot.ID,
		State: model.RaybotStateOffline,
	})

	if err != nil {
		c.logger.Error("Error updating raybot status", slog.Any("error", err))
	}

	close(c.send)
}

func (c *RaybotClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.logger.Error("Error when write message", slog.Any("error", err))
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				c.logger.Warn("Error when ping", slog.Any("error", err))
				return
			}
		}
	}
}

func (c *RaybotClient) ReadPump(unregister chan<- *RaybotClient) {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("Unexpected close", slog.Any("error", err))
			}
			unregister <- c
			break
		}
		handleMessage(c, message)
	}
}

func (c *RaybotClient) SendCommand(cmd *model.RaybotCommand) error {
	c.logger.Debug("Received command", slog.Any("cmd", cmd))
	c.cmdMap[cmd.ID.String()] = cmd

	msg := OutboundCommandMsg{
		ID:   cmd.ID.String(),
		Type: cmd.Type,
		Data: map[string]interface{}{
			// TODO: Need to fix this more type safe
			"damn": "damn",
			"test": "test",
		},
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.send <- msgJSON
	return nil
}

func handleMessage(c *RaybotClient, msg []byte) {
	// Identify inbound message
	var inboundMsg InboundMsg
	err := json.Unmarshal(msg, &inboundMsg)
	if err != nil {
		closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid message")
		return
	}

	// Route message based on operation
	switch inboundMsg.Operation {
	case OperationPublish:
		handlePublish(c, inboundMsg)
	case OperationResponse:
		// handleResponse(c, inboundMsg)
	default:
		closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid operation")
		return
	}
}

func closeConn(c *RaybotClient, statusCode int, text string) {
	err := c.conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))

	if err != nil {
		c.logger.Error("Error when close connection", slog.Any("error", err))
	}
}
