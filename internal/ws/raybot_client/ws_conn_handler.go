package raybotclient

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

const (
	inboundMsgWorkerCount = 2
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

func (ws *WebSocket) HandleConnection(w http.ResponseWriter, r *http.Request) {
	raybotID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		ws.log.Info("Invalid raybot ID", slog.Any("error", err))
		return
	}

	c := ws.getClient(raybotID)
	if c != nil {
		ws.log.Info("Raybot is already connected", slog.String("raybot_id", raybotID.String()))
		return
	}

	// Get raybot
	raybotModel, err := ws.raybotRepo.Get(r.Context(), raybotID)
	if err != nil {
		if xerrors.IsNotFound(err) {
			ws.log.Info("Raybot not found", slog.String("raybot_id", raybotID.String()))
			return
		}
		ws.log.Error("Failed to get raybot", slog.Any("error", err))
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.log.Error("Failed to upgrade connection", slog.Any("error", err))
		return
	}

	defer func() {
		ws.log.Info("Raybot disconnected", slog.String("raybot_id", raybotID.String()))
		raybotModel.Status = model.RaybotStatusOffline
		raybotModel, err = ws.raybotRepo.Update(r.Context(), raybotModel)
		if err != nil {
			ws.log.Error("Failed to update raybot status", slog.Any("error", err))
		}

		r.Context().Done()
		conn.Close()
	}()

	// Update raybot info
	raybotModel.Status = model.RaybotStatusIdle
	ip := getClientIPAddress(r)
	raybotModel.IpAddress = &ip
	now := time.Now()
	raybotModel.LastConnectedAt = &now
	raybotModel, err = ws.raybotRepo.Update(r.Context(), raybotModel)
	if err != nil {
		ws.log.Error("Failed to update raybot status", slog.Any("error", err))
		return
	}

	// Create raybot client
	c = NewRaybotClient(&raybotModel, conn, ws.raybotCommandSvc, ws.log)
	ws.addClient(c)
	defer ws.removeClient(raybotID)

	// Start worker
	StartInboundMsgWorker(r.Context(), c, inboundMsgWorkerCount)
	go c.WritePump()
	c.ReadPump() // Block until connection is closed
}
