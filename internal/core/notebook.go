package core

import "time"

// NotebookID is a unique identifier for a Notebook.
type NotebookID string

// String returns the string representation of the NotebookID.
func (id NotebookID) String() string {
	return string(id)
}

// Notebook represents a notebook that groups cards for study purposes.
type Notebook struct {
	id        NotebookID
	name      string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

// NewNotebook creates a new Notebook with the given id, name and creator.
// It initializes createdAt and updatedAt to the current time.
func NewNotebook(id NotebookID, name string, creator UserID) *Notebook {
	now := time.Now()
	return &Notebook{
		id:        id,
		name:      name,
		creator:   creator,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the notebook.
func (n *Notebook) ID() NotebookID {
	return n.id
}

// Name returns the notebook's display name.
func (n *Notebook) Name() string {
	return n.name
}

// Creator returns the UserID of the user who created this notebook.
func (n *Notebook) Creator() UserID {
	return n.creator
}

// CreatedAt returns the timestamp when the notebook was created.
func (n *Notebook) CreatedAt() time.Time {
	return n.createdAt
}

// UpdatedAt returns the timestamp when the notebook was last modified.
func (n *Notebook) UpdatedAt() time.Time {
	return n.updatedAt
}

// SetName updates the notebook's name and sets updatedAt to the current time.
func (n *Notebook) SetName(name string) {
	n.name = name
	n.updatedAt = time.Now()
}
