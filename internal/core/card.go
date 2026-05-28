package core

import "time"

// CardID is a unique identifier for a Card.
type CardID string

// String returns the string representation of the CardID.
func (id CardID) String() string {
	return string(id)
}

// Card represents a flashcard with a hint and content.
type Card struct {
	id        CardID
	hint      string
	content   string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

// NewCard creates a new Card with the given id, hint, content and creator.
// It initializes createdAt and updatedAt to the current time.
func NewCard(id CardID, hint string, content string, creator UserID) *Card {
	now := time.Now()
	return &Card{
		id:        id,
		hint:      hint,
		content:   content,
		creator:   creator,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the card.
func (c *Card) ID() CardID {
	return c.id
}

// Hint returns the card's hint or prompt.
func (c *Card) Hint() string {
	return c.hint
}

// Content returns the card's main content or answer.
func (c *Card) Content() string {
	return c.content
}

// Creator returns the UserID of the user who created this card.
func (c *Card) Creator() UserID {
	return c.creator
}

// CreatedAt returns the timestamp when the card was created.
func (c *Card) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the timestamp when the card was last modified.
func (c *Card) UpdatedAt() time.Time {
	return c.updatedAt
}

// SetHint updates the card's hint and sets updatedAt to the current time.
func (c *Card) SetHint(hint string) {
	c.hint = hint
	c.updatedAt = time.Now()
}

// SetContent updates the card's content and sets updatedAt to the current time.
func (c *Card) SetContent(content string) {
	c.content = content
	c.updatedAt = time.Now()
}
