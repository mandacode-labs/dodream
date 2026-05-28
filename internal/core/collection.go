package core

import (
	"time"

	"github.com/google/uuid"
)

// CollectionID is a unique identifier for a Collection.
type CollectionID string

// String returns the string representation of the CollectionID.
func (id CollectionID) String() string {
	return string(id)
}

// Collection represents a personalized grouping of cards for study purposes.
type Collection struct {
	id        CollectionID
	name      string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

// NewCollection creates a new Collection with the given name and creator.
// It auto-generates a UUID for the id and initializes createdAt and updatedAt.
func NewCollection(name string, creator UserID) *Collection {
	now := time.Now()
	return &Collection{
		id:        CollectionID(uuid.New().String()),
		name:      name,
		creator:   creator,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the collection.
func (c *Collection) ID() CollectionID {
	return c.id
}

// Name returns the collection's display name.
func (c *Collection) Name() string {
	return c.name
}

// Creator returns the UserID of the user who created this collection.
func (c *Collection) Creator() UserID {
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

// SetName updates the collection's name and sets updatedAt to the current time.
func (c *Collection) SetName(name string) {
	c.name = name
	c.updatedAt = time.Now()
}
