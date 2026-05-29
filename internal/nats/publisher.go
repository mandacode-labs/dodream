// Package nats provides generic NATS messaging infrastructure.
package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	// SubjectStudyEvents is the NATS subject for study events.
	SubjectStudyEvents = "study.events"
)

// StudyEventMessage is the NATS message payload for study events.
type StudyEventMessage struct {
	EventID        string  `json:"event_id"`
	CollectionID   string  `json:"collection_id"`
	CardID         string  `json:"card_id"`
	UserID         string  `json:"user_id"`
	EventType      string  `json:"event_type"`
	Quality        *int    `json:"quality,omitempty"`
	ResponseTimeNs *int64  `json:"response_time_ns,omitempty"`
	PreviousDeckID *string `json:"previous_deck_id,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

// Publisher handles publishing messages to NATS.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher creates a new NATS publisher.
func NewPublisher(natsURL string) (*Publisher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	return &Publisher{conn: conn}, nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	p.conn.Close()
}

// PublishStudyEvent publishes a study event message to NATS.
func (p *Publisher) PublishStudyEvent(ctx context.Context, msg *StudyEventMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal study event: %w", err)
	}

	if err := p.conn.Publish(SubjectStudyEvents, data); err != nil {
		return fmt.Errorf("failed to publish study event: %w", err)
	}

	return nil
}
