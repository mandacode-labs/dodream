package core

import "time"

// UserID is a unique identifier for a User.
type UserID string

// String returns the string representation of the UserID.
func (id UserID) String() string {
	return string(id)
}

// User represents a registered user in the system.
type User struct {
	id         UserID
	nickname   string
	providerID string
	createdAt  time.Time
	updatedAt  time.Time
}

// NewUser creates a new User with the given id, nickname and providerID.
// It initializes createdAt and updatedAt to the current time.
func NewUser(id UserID, nickname string, providerID string) *User {
	now := time.Now()
	return &User{
		id:         id,
		nickname:   nickname,
		providerID: providerID,
		createdAt:  now,
		updatedAt:  now,
	}
}

// ID returns the unique identifier of the user.
func (u *User) ID() UserID {
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
