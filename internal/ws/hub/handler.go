package hub

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	commandSvc "github.com/tuanvumaihuynh/roboflow/internal/command/service"
	raybotSvc "github.com/tuanvumaihuynh/roboflow/internal/raybot/service"
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
	raybotSvc raybotSvc.RaybotService,
	cmdSvc commandSvc.CommandService,
) {
	ctx := r.Context()
	raybotID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		h.logger.Error("Error parsing raybot id", zap.Error(err))
		return
	}

	// Check if raybot is already connected
	_, exist := h.clients[raybotID.String()]
	if exist {
		h.logger.Error("Raybot client is already connected", zap.String("id", raybotID.String()))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Error upgrading connection", zap.Error(err))
		return
	}

	// Create a new RaybotClient
	client, err := client.InitilizeRaybotClient(ctx, raybotID, raybotSvc, cmdSvc, conn, h.logger)
	if err != nil {
		h.logger.Error("Error initializing raybot client", zap.Error(err))
		return
	}
	h.register <- client
	h.logger.Info("Raybot connected", zap.String("id", raybotID.String()))

	go client.WritePump()
	go client.ReadPump(h.unregister)
}
