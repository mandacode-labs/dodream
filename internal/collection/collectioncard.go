// Package collection provides collection membership for cards.
package collection

import (
	"time"

	"github.com/google/uuid"
)

// CollectionCardID is a unique identifier for a CollectionCard.
type CollectionCardID string

// String returns the string representation of the ID.
func (id CollectionCardID) String() string {
	return string(id)
}

// CollectionCard represents a card within a collection with an optional deck reference.
type CollectionCard struct {
	id           CollectionCardID
	collectionID ID
	cardID       string // card ID as string to avoid import cycle
	deckID       *string
	createdAt    time.Time
	updatedAt    time.Time
}

// NewCollectionCard creates a new CollectionCard.
func NewCollectionCard(collectionID ID, cardID string, deckID *string) *CollectionCard {
	now := time.Now()
	return &CollectionCard{
		id:           CollectionCardID(uuid.New().String()),
		collectionID: collectionID,
		cardID:       cardID,
		deckID:       deckID,
		createdAt:    now,
		updatedAt:    now,
	}
}

// ID returns the unique identifier of the collection card.
func (cc *CollectionCard) ID() CollectionCardID {
	return cc.id
}

// CollectionID returns the collection identifier.
func (cc *CollectionCard) CollectionID() ID {
	return cc.collectionID
}

// CardID returns the card identifier.
func (cc *CollectionCard) CardID() string {
	return cc.cardID
}

// DeckID returns the optional deck identifier.
func (cc *CollectionCard) DeckID() *string {
	return cc.deckID
}

// CreatedAt returns the timestamp when the collection card was created.
func (cc *CollectionCard) CreatedAt() time.Time {
	return cc.createdAt
}

// UpdatedAt returns the timestamp when the collection card was last modified.
func (cc *CollectionCard) UpdatedAt() time.Time {
	return cc.updatedAt
}

// SetDeckID updates the deck reference.
func (cc *CollectionCard) SetDeckID(deckID *string) {
	cc.deckID = deckID
	cc.updatedAt = time.Now()
}

// NewCollectionCardWithID creates a CollectionCard from database values.
func NewCollectionCardWithID(id CollectionCardID, collectionID ID, cardID string, deckID *string, createdAt, updatedAt time.Time) *CollectionCard {
	return &CollectionCard{
		id:           id,
		collectionID: collectionID,
		cardID:       cardID,
		deckID:       deckID,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
