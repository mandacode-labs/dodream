package core

import "time"

type NotebookID string

func (id NotebookID) String() string {
	return string(id)
}

type Notebook struct {
	id        NotebookID
	name      string
	creator   UserID
	createdAt time.Time
	updatedAt time.Time
}

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

func (n *Notebook) ID() NotebookID {
	return n.id
}

func (n *Notebook) Name() string {
	return n.name
}

func (n *Notebook) Creator() UserID {
	return n.creator
}

func (n *Notebook) CreatedAt() time.Time {
	return n.createdAt
}

func (n *Notebook) UpdatedAt() time.Time {
	return n.updatedAt
}

func (n *Notebook) SetName(name string) {
	n.name = name
	n.updatedAt = time.Now()
}
