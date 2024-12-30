package raybotclient

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
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

// RaybotClient represents a raybot client connect through websocket.
type RaybotClient struct {
	raybot *model.Raybot

	// The current command that the raybot is processing.
	// Now RaybotClient can processing only one command at a time.
	currentCmd   *model.RaybotCommand
	currentCmdMu sync.Mutex

	// The websocket connection.
	conn *websocket.Conn

	// Services
	raybotCommandSvc raybotcommand.Service

	// Buffered channel of unprocessed inbound messages.
	inboundChan chan []byte

	// Buffered channel of outbound messages.
	outboundChan chan []byte

	// Logger
	log *slog.Logger
}

func NewRaybotClient(
	raybot *model.Raybot,
	conn *websocket.Conn,
	raybotCommandSvc raybotcommand.Service,
	log *slog.Logger,
) *RaybotClient {
	return &RaybotClient{
		raybot:           raybot,
		conn:             conn,
		raybotCommandSvc: raybotCommandSvc,
		inboundChan:      make(chan []byte, 256),
		outboundChan:     make(chan []byte, 256),
		log:              log,
	}
}

func (c *RaybotClient) GetCurrentCmd() *model.RaybotCommand {
	c.currentCmdMu.Lock()
	defer c.currentCmdMu.Unlock()

	return c.currentCmd
}

func (c *RaybotClient) SetCurrentCmd(cmd *model.RaybotCommand) {
	c.currentCmdMu.Lock()
	defer c.currentCmdMu.Unlock()

	c.currentCmd = cmd
}

func (c *RaybotClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.outboundChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.log.Error("Error when write message", slog.Any("error", err))
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

func (c *RaybotClient) ReadPump() {
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
				c.log.Error("Unexpected close", slog.Any("error", err))
			}
			break
		}

		select {
		case c.inboundChan <- message:
		default:
			c.log.Error("Inbound channel is full, drop message",
				slog.Int("inbound_chan_size", len(c.inboundChan)))
		}
	}
}

func closeConn(c *RaybotClient, statusCode int, text string) {
	err := c.conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))

	if err != nil {
		c.log.Error("Error when close connection", slog.Any("error", err))
	}
}
