// Package identity provides strongly-typed identifiers to prevent mixing
// different entity IDs at compile time.
//
// Usage: prefer these types over raw strings when passing IDs across domain
// boundaries (e.g. creatorID fields, cross-domain references).
package identity

// UserID is a typed identifier for users.
type UserID string

// String returns the string representation of the ID.
func (id UserID) String() string { return string(id) }

// CardID is a typed identifier for cards.
type CardID string

// String returns the string representation of the ID.
func (id CardID) String() string { return string(id) }

// DeckID is a typed identifier for decks.
type DeckID string

// String returns the string representation of the ID.
func (id DeckID) String() string { return string(id) }

// CollectionID is a typed identifier for collections.
type CollectionID string

// String returns the string representation of the ID.
func (id CollectionID) String() string { return string(id) }

// StateID is a typed identifier for states.
type StateID string

// String returns the string representation of the ID.
func (id StateID) String() string { return string(id) }

// StudyEventID is a typed identifier for study events.
type StudyEventID string

// String returns the string representation of the ID.
func (id StudyEventID) String() string { return string(id) }
