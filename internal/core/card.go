package core

import "time"

type CardID string

func (id CardID) String() string {
	return string(id)
}

type Card struct {
	id        CardID
	hint      string
	content   string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

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

func (c *Card) ID() CardID {
	return c.id
}

func (c *Card) Hint() string {
	return c.hint
}

func (c *Card) Content() string {
	return c.content
}

func (c *Card) Creator() UserID {
	return c.creator
}

func (c *Card) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Card) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Card) SetHint(hint string) {
	c.hint = hint
	c.updatedAt = time.Now()
}

func (c *Card) SetContent(content string) {
	c.content = content
	c.updatedAt = time.Now()
}
