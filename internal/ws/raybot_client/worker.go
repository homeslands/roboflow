package raybotclient

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
)

// StartInboundMsgWorker starts the inbound message worker.
// Non-blocking function.
func StartInboundMsgWorker(ctx context.Context, c *RaybotClient, numWorkers uint8) {
	for i := uint8(0); i < numWorkers; i++ {
		go func(workerID uint8) {
			c.log.Info("Inbound msg worker started", slog.Int("worker_id", int(workerID)))

			for {
				select {
				case <-ctx.Done():
					c.log.Info("Worker stopped", slog.Int("worker_id", int(workerID)))
					return
				case msg, ok := <-c.inboundChan:
					if !ok {
						// Channel closed, stop the worker
						c.log.Info("Inbound channel closed", slog.Int("worker_id", int(workerID)))
						return
					}
					processInboundMsg(c, msg)
				}
			}
		}(i)
	}
}

func processInboundMsg(c *RaybotClient, data []byte) {
	var temp struct {
		Op Operation `json:"op"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		c.log.Error("Error unmarshal message", slog.Any("error", err))
		closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid message")
		return
	}

	switch temp.Op {
	case OperationPublish:
		var msg InboundPublishMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			c.log.Error("Error unmarshal publish message", slog.Any("error", err))
			closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid publish message")
			return
		}
		handleInboundPublishMsg(c, msg)

	case OperationResponse:
		var msg InboundResponseMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			c.log.Error("Error unmarshal response message", slog.Any("error", err))
			closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid response message")
			return
		}
		handleInboundResponseMsg(c, msg)

	default:
		c.log.Error("Invalid operation", slog.String("op", string(temp.Op)))
		closeConn(c, websocket.CloseInvalidFramePayloadData, "Invalid operation")
	}
}
