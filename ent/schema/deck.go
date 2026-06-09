package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Deck holds the schema definition for the Deck entity.
type Deck struct {
	ent.Schema
}

// Fields of the Deck.
func (Deck) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Deck.
func (Deck) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).
			Ref("decks").
			Unique().
			Required(),
		edge.From("cards", Card.Type).
			Ref("decks"),
		edge.To("collection_cards", CollectionCard.Type),
		edge.To("study_events", StudyEvent.Type),
	}
}
