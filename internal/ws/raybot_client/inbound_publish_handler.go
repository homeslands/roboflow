package raybotclient

import "log/slog"

func handleInboundPublishMsg(c *RaybotClient, msg InboundPublishMsg) {
	c.log.Debug("Received publish msg",
		slog.String("topic", string(msg.Topic)),
		slog.String("data", string(msg.Data)),
	)

	// TODO: Implement based on requirements
	c.log.Warn("OperationPublish is not implemented")
}
