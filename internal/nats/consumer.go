// Package nats provides generic NATS messaging infrastructure.
package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// Handler is the function signature for processing study events.
type Handler func(msg *StudyEventMessage) error

// Consumer handles consuming study events from NATS.
type Consumer struct {
	conn       *nats.Conn
	sub        *nats.Subscription
	handler    Handler
	queueGroup string
}

// NewConsumer creates a new NATS consumer.
func NewConsumer(natsURL string, queueGroup string, handler Handler) (*Consumer, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	return &Consumer{
		conn:       conn,
		handler:    handler,
		queueGroup: queueGroup,
	}, nil
}

// Start begins consuming messages from NATS.
func (c *Consumer) Start() error {
	sub, err := c.conn.QueueSubscribe(SubjectStudyEvents, c.queueGroup, func(msg *nats.Msg) {
		var event StudyEventMessage
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("nats: failed to unmarshal message: %v", err)
			return
		}

		if err := c.handler(&event); err != nil {
			log.Printf("nats: handler error for event %s: %v", event.EventID, err)
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	c.sub = sub
	return nil
}

// Stop unsubscribes and closes the connection.
func (c *Consumer) Stop() error {
	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			return err
		}
	}
	c.conn.Close()
	return nil
}
