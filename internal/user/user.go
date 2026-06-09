// Package user provides the user domain.
package user

import (
	"time"

	"github.com/google/uuid"
)

// ID is a unique identifier for a User.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// User represents a registered user in the system.
type User struct {
	id         ID
	nickname   string
	providerID string
	createdAt  time.Time
	updatedAt  time.Time
}

// New creates a new User with the given nickname and providerID.
func New(nickname string, providerID string) *User {
	now := time.Now()
	return &User{
		id:         ID(uuid.New().String()),
		nickname:   nickname,
		providerID: providerID,
		createdAt:  now,
		updatedAt:  now,
	}
}

// ID returns the unique identifier of the user.
func (u *User) ID() ID {
	return u.id
}

// Nickname returns the user's display name.
func (u *User) Nickname() string {
	return u.nickname
}

// ProviderID returns the external provider identifier for the user.
func (u *User) ProviderID() string {
	return u.providerID
}

// CreatedAt returns the timestamp when the user was created.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the timestamp when the user was last modified.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// SetNickname updates the user's nickname and sets updatedAt to the current time.
func (u *User) SetNickname(nickname string) {
	u.nickname = nickname
	u.updatedAt = time.Now()
}

// SetProviderID updates the user's providerID and sets updatedAt to the current time.
func (u *User) SetProviderID(providerID string) {
	u.providerID = providerID
	u.updatedAt = time.Now()
}

// NewWithID creates a User from database values.
func NewWithID(id ID, nickname string, providerID string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:         id,
		nickname:   nickname,
		providerID: providerID,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}
