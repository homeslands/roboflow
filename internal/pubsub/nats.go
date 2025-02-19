package pubsub

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	ns "github.com/nats-io/nats-server/v2/server"
	nc "github.com/nats-io/nats.go"

	"github.com/tuanvumaihuynh/roboflow/pkg/config"
)

// NewNatsPubSub creates a new nats publisher and subscriber
func NewNatsPubSub(conf config.NatsConfig, log *slog.Logger) (*nats.Publisher, *nats.Subscriber, error) {
	var server *ns.Server
	var err error
	url := nc.DefaultURL

	clientOpts := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.ReconnectWait(1 * time.Second),
	}

	// If the URL is not provided, we will create a new embedded server.
	if conf.URL == nil {
		server, err = newEmbeddedNatsServer()
		if err != nil {
			return nil, nil, fmt.Errorf("create nats server: %w", err)
		}

		if conf.EnableLog {
			server.ConfigureLogger()
		}

		clientOpts = append(clientOpts, nc.InProcessServer(server))

		log.Info("Starting NATS server in embedded mode")
		go server.Start()
		if !server.ReadyForConnections(5 * time.Second) {
			return nil, nil, errors.New("nats server not ready for connections")
		}
	} else {
		url = *conf.URL
	}

	wLog := watermill.NewSlogLogger(log)
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}
	jetstreamConf := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       true,
		AckAsync:         true,
		DurablePrefix:    "roboflow",
	}
	subscriber, err := nats.NewSubscriber(nats.SubscriberConfig{
		URL:            url,
		CloseTimeout:   30 * time.Second,
		AckWaitTimeout: 30 * time.Second,
		NatsOptions:    clientOpts,
		Unmarshaler:    nats.GobMarshaler{},
		JetStream:      jetstreamConf,
	}, wLog)
	if err != nil {
		return nil, nil, fmt.Errorf("create nats subscriber: %w", err)
	}

	publisher, err := nats.NewPublisher(nats.PublisherConfig{
		URL:         url,
		NatsOptions: clientOpts,
		Marshaler:   nats.GobMarshaler{},
		JetStream:   jetstreamConf,
	}, wLog)
	if err != nil {
		return nil, nil, fmt.Errorf("create nats publisher: %w", err)
	}

	return publisher, subscriber, nil
}

func newEmbeddedNatsServer() (*ns.Server, error) {
	opts := &ns.Options{
		ServerName:      "roboflow-nats",
		DontListen:      true,
		JetStream:       true,
		JetStreamDomain: "roboflow-nats-embedded",
	}

	server, err := ns.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create nats server: %w", err)
	}

	return server, nil
}
