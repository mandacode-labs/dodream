package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CollectionCard holds the schema definition for the CollectionCard entity.
type CollectionCard struct {
	ent.Schema
}

// Fields of the CollectionCard.
func (CollectionCard) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the CollectionCard.
func (CollectionCard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("collection", Collection.Type).
			Ref("collection_cards").
			Unique().
			Required(),
		edge.From("card", Card.Type).
			Ref("collection_cards").
			Unique().
			Required(),
		edge.From("deck", Deck.Type).
			Ref("collection_cards").
			Unique(),
	}
}
