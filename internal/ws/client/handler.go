package client

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func handlePublish(c *RaybotClient, msg InboundMsg) {
	// Route message based on topic
	topic := *msg.Topic
	switch topic {
	// case topicStatus:
	// 	handleStatus(c, msg)
	// case topicRamCpu:
	// 	handleRamCpu(c, msg)
	// case topicLog:
	// 	handleLog(c, msg)
	// case topicBatterySensor:
	// 	handleBatterySensor(c, msg)
	// case topicWeightSensor:
	// 	handleWeightSensor(c, msg)
	// case topicForwardDistanceSensor:
	// 	handleForwardDistanceSensor(c, msg)
	// case topicBackwardDistanceSensor:
	// 	handleBackwardDistanceSensor(c, msg)
	// case topicMovementMotor:
	// 	handleMovementMotor(c, msg)
	// case topicLiftMotor:
	// 	handleLiftMotor(c, msg)
	default:
		closeConn(c, websocket.CloseUnsupportedData, "unsupported topic")
	}
}

func handleResponse(c *RaybotClient, msg InboundMsg) {
	c.logger.Debug("Received msg", zap.Any("msg", msg))
	ID, err := uuid.Parse(*msg.Id)
	if err != nil {
		c.logger.Error("Error parsing id", zap.Error(err))
		closeConn(c, websocket.CloseInvalidFramePayloadData, "failed to parse id")
		return
	}

	cmd, err := c.cmdSvc.GetCommand(context.Background(), ID)
	if err != nil {
		if xerrors.IsNotFound(err) {
			c.logger.Error("Command not found", zap.Error(err))
			closeConn(c, websocket.CloseInvalidFramePayloadData, "command not found")
			return
		}
		c.logger.Error("Error getting command", zap.Error(err))
		closeConn(c, websocket.CloseInternalServerErr, "failed to get command")
		return
	}

	// Get command status
	var partialData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &partialData); err != nil {
		panic(err)
	}

	cmd.CompletedAt = nil
	// // Get status
	// var status command.Status
	// if stat, exists := partialData["status"].(string); exists {
	// 	status = command.Status(stat)
	// }

	// // Set status to command
	// cmd.Status = status

	// // Route message based on command type
	// switch cmd.Type {
	// case command.CommandTypeMoveForward:
	// 	handleCommandMoveForward(c, *cmd)
	// case command.CommandTypeMoveBackward:
	// 	handleCommandMoveBackward(c, *cmd)
	// default:
	// 	closeConn(c, websocket.CloseUnsupportedData, "unsupported command type")
	// }
}
