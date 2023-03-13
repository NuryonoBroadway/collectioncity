package collection_core

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
)

type Message struct {
	Topic     string
	Data      []byte
	Attribute map[string]string
}

type gPublisher struct {
	client *pubsub.Client
}

// NewGPublisher create new instance of pubsub publisher
func NewGPublisher(client *pubsub.Client) *gPublisher {
	return &gPublisher{client: client}
}

// Publish publish message to the topic
func (p *gPublisher) Publish(ctx context.Context, msg *Message) error {
	tp := p.client.Topic(msg.Topic)

	payload := &pubsub.Message{
		Data:        msg.Data,
		Attributes:  msg.Attribute,
		PublishTime: time.Now(),
	}

	_, err := tp.Publish(ctx, payload).Get(ctx)
	return err
}
