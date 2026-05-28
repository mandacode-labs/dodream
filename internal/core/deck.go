package core

import "time"

type DeckID string

func (id DeckID) String() string {
	return string(id)
}

type Deck struct {
	id        DeckID
	name      string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

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

func (d *Deck) ID() DeckID {
	return d.id
}

func (d *Deck) Name() string {
	return d.name
}

func (d *Deck) Creator() UserID {
	return d.creator
}

func (d *Deck) CreatedAt() time.Time {
	return d.createdAt
}

func (d *Deck) UpdatedAt() time.Time {
	return d.updatedAt
}

func (d *Deck) SetName(name string) {
	d.name = name
	d.updatedAt = time.Now()
}
