package hub

import (
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/internal/ws/client"
)

type Hub struct {
	clients    map[string]*client.RaybotClient
	register   chan *client.RaybotClient
	unregister chan *client.RaybotClient

	logger *zap.Logger
}

type HubConfig struct {
	Logger *zap.Logger
}

func NewHub(cfg HubConfig) *Hub {
	h := &Hub{
		clients:    make(map[string]*client.RaybotClient),
		register:   make(chan *client.RaybotClient),
		unregister: make(chan *client.RaybotClient),
		logger:     cfg.Logger.With(zap.String("module", "WsHub")),
	}

	go h.MainLoop()
	return h
}

func (h *Hub) MainLoop() {
	for {
		select {
		case c := <-h.register:
			id := c.ID()
			h.clients[id] = c
			h.logger.Info("Raybot client registered", zap.String("client_id", id))
		case c := <-h.unregister:
			id := c.ID()
			if _, ok := h.clients[id]; ok {
				delete(h.clients, id)
				c.Close()
				h.logger.Info("Raybot client unregistered", zap.String("client_id", id))
			}
		}
	}
}
