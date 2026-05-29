// Package deck provides the deck domain.
package deck

import (
	"time"

	"github.com/google/uuid"
)

// ID is a unique identifier for a Deck.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// Deck represents a collection of flashcards.
type Deck struct {
	id        ID
	name      string
	creatorID string // user ID as string to avoid import cycle
	createdAt time.Time
	updatedAt time.Time
}

// New creates a new Deck with the given name and creator.
func New(name string, creatorID string) *Deck {
	now := time.Now()
	return &Deck{
		id:        ID(uuid.New().String()),
		name:      name,
		creatorID: creatorID,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the deck.
func (d *Deck) ID() ID {
	return d.id
}

// Name returns the deck's display name.
func (d *Deck) Name() string {
	return d.name
}

// CreatorID returns the user ID who created this deck.
func (d *Deck) CreatorID() string {
	return d.creatorID
}

// CreatedAt returns the timestamp when the deck was created.
func (d *Deck) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the timestamp when the deck was last modified.
func (d *Deck) UpdatedAt() time.Time {
	return d.updatedAt
}

// SetName updates the deck's name.
func (d *Deck) SetName(name string) {
	d.name = name
	d.updatedAt = time.Now()
}

// NewWithID creates a Deck from database values.
func NewWithID(id ID, name string, creatorID string, createdAt, updatedAt time.Time) *Deck {
	return &Deck{
		id:        id,
		name:      name,
		creatorID: creatorID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
