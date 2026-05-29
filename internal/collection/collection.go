// Package collection provides the collection domain.
package collection

import (
	"time"

	"github.com/google/uuid"
)

// ID is a unique identifier for a Collection.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// Collection represents a personalized grouping of cards for study purposes.
type Collection struct {
	id        ID
	name      string
	creator   string // user ID as string to avoid import cycle
	createdAt time.Time
	updatedAt time.Time
}

// New creates a new Collection with the given name and creator.
func New(name string, creator string) *Collection {
	now := time.Now()
	return &Collection{
		id:        ID(uuid.New().String()),
		name:      name,
		creator:   creator,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the collection.
func (c *Collection) ID() ID {
	return c.id
}

// Name returns the collection's display name.
func (c *Collection) Name() string {
	return c.name
}

// Creator returns the user ID who created this collection.
func (c *Collection) Creator() string {
	return c.creator
}

// CreatedAt returns the timestamp when the collection was created.
func (c *Collection) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the timestamp when the collection was last modified.
func (c *Collection) UpdatedAt() time.Time {
	return c.updatedAt
}

// SetName updates the collection's name.
func (c *Collection) SetName(name string) {
	c.name = name
	c.updatedAt = time.Now()
}

// NewWithID creates a Collection from database values.
func NewWithID(id ID, name string, creator string, createdAt, updatedAt time.Time) *Collection {
	return &Collection{
		id:        id,
		name:      name,
		creator:   creator,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
