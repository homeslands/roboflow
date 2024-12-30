package raybotclient

import (
	"log/slog"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
)

type WebSocket struct {
	// Map of connected raybot clients.
	clients map[uuid.UUID]*RaybotClient
	mu      sync.RWMutex

	raybotRepo       model.RaybotRepository
	raybotCommandSvc raybotcommand.Service
	log              *slog.Logger
}

func NewWebSocket(
	raybotRepo model.RaybotRepository,
	raybotCommandSvc raybotcommand.Service,
	log *slog.Logger,
) *WebSocket {
	return &WebSocket{
		clients:          make(map[uuid.UUID]*RaybotClient),
		raybotRepo:       raybotRepo,
		raybotCommandSvc: raybotCommandSvc,
		log:              log,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws-raybot", ws.HandleConnection)
}

func (ws *WebSocket) getClient(raybotID uuid.UUID) *RaybotClient {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	return ws.clients[raybotID]
}

func (ws *WebSocket) addClient(client *RaybotClient) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.clients[client.raybot.ID] = client
}

func (ws *WebSocket) removeClient(raybotID uuid.UUID) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	delete(ws.clients, raybotID)
}
