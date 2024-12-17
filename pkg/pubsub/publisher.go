package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(topic string, payload []byte) error
}

type eventPublisher struct {
	publisher message.Publisher
}

// NewPublisher creates a new event publisher
func NewPublisher(p message.Publisher) *eventPublisher {
	if p == nil {
		panic("nil publisher")
	}
	return &eventPublisher{
		publisher: p,
	}
}

func (p *eventPublisher) Publish(topic string, payload []byte) error {
	msg := message.NewMessage(uuid.NewString(), payload)
	return p.publisher.Publish(topic, msg)
}
