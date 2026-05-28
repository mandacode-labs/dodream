package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NotebookCard holds the schema definition for the NotebookCard entity.
type NotebookCard struct {
	ent.Schema
}

// Fields of the NotebookCard.
func (NotebookCard) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the NotebookCard.
func (NotebookCard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("notebook", Notebook.Type).
			Ref("notebook_cards").
			Unique().
			Required(),
		edge.From("card", Card.Type).
			Ref("notebook_cards").
			Unique().
			Required(),
		edge.From("deck", Deck.Type).
			Ref("notebook_cards").
			Unique(),
	}
}
