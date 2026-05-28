package core

import "time"

// DeckID is a unique identifier for a Deck.
type DeckID string

// String returns the string representation of the DeckID.
func (id DeckID) String() string {
	return string(id)
}

// Deck represents a collection of flashcards.
type Deck struct {
	id        DeckID
	name      string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

// NewDeck creates a new Deck with the given id, name and creator.
// It initializes createdAt and updatedAt to the current time.
func NewDeck(id DeckID, name string, creator UserID) *Deck {
	now := time.Now()
	return &Deck{
		id:        id,
		name:      name,
		creator:   creator,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the deck.
func (d *Deck) ID() DeckID {
	return d.id
}

// Name returns the deck's display name.
func (d *Deck) Name() string {
	return d.name
}

// Creator returns the UserID of the user who created this deck.
func (d *Deck) Creator() UserID {
	return d.creator
}

// CreatedAt returns the timestamp when the deck was created.
func (d *Deck) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the timestamp when the deck was last modified.
func (d *Deck) UpdatedAt() time.Time {
	return d.updatedAt
}

// SetName updates the deck's name and sets updatedAt to the current time.
func (d *Deck) SetName(name string) {
	d.name = name
	d.updatedAt = time.Now()
}
