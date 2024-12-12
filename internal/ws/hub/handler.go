package hub

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	raybotSvc "github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	raybotCommandSvc "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/ws/client"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}
)

func (h *Hub) HandleRaybotClientConnect(
	w http.ResponseWriter,
	r *http.Request,
	raybotSvc raybotSvc.Service,
	raybotCommanddSvc raybotCommandSvc.Service,
) {
	ctx := r.Context()
	raybotID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		h.logger.Error("Error parsing raybot id", slog.Any("error", err.Error()))
		return
	}

	// Check if raybot is already connected
	_, exist := h.clients[raybotID.String()]
	if exist {
		h.logger.Error("Raybot client is already connected", slog.String("id", raybotID.String()))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Error upgrading connection", slog.Any("error", err))
		return
	}

	// Create a new RaybotClient
	client, err := client.InitilizeRaybotClient(ctx, raybotID, raybotSvc, raybotCommanddSvc, conn, h.logger)
	if err != nil {
		h.logger.Error("Error initializing raybot client", slog.Any("error", err))
		return
	}
	h.register <- client
	h.logger.Info("Raybot connected", slog.String("id", raybotID.String()))

	go client.WritePump()
	go client.ReadPump(h.unregister)
}
