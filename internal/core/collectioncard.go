package core

import "time"

// CollectionCardID is a unique identifier for a CollectionCard.
type CollectionCardID string

// String returns the string representation of the CollectionCardID.
func (id CollectionCardID) String() string {
	return string(id)
}

// CollectionCard represents a card within a collection with an optional deck reference.
type CollectionCard struct {
	id           CollectionCardID
	collectionID CollectionID
	cardID       CardID
	deckID       *DeckID
	createdAt    time.Time
	updatedAt    time.Time
}

// NewCollectionCard creates a new CollectionCard with the given id, collectionID, cardID, and optional deckID.
// It initializes createdAt and updatedAt to the current time.
func NewCollectionCard(
	id CollectionCardID,
	collectionID CollectionID,
	cardID CardID,
	deckID *DeckID,
) *CollectionCard {
	now := time.Now()
	return &CollectionCard{
		id:           id,
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
func (cc *CollectionCard) CollectionID() CollectionID {
	return cc.collectionID
}

// CardID returns the card identifier.
func (cc *CollectionCard) CardID() CardID {
	return cc.cardID
}

// DeckID returns the optional deck identifier.
func (cc *CollectionCard) DeckID() *DeckID {
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

// SetDeckID updates the deck reference and sets updatedAt to the current time.
func (cc *CollectionCard) SetDeckID(deckID *DeckID) {
	cc.deckID = deckID
	cc.updatedAt = time.Now()
}
