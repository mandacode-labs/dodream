package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StudyEvent holds the schema definition for the StudyEvent entity.
type StudyEvent struct {
	ent.Schema
}

// Fields of the StudyEvent.
func (StudyEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("event_type").NotEmpty(),
		field.Int("quality").Optional().Nillable(),
		field.Int64("response_time_ns").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the StudyEvent.
func (StudyEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("collection", Collection.Type).
			Ref("study_events").
			Unique().
			Required(),
		edge.From("card", Card.Type).
			Ref("study_events").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("study_events").
			Unique().
			Required(),
		edge.From("previous_deck", Deck.Type).
			Ref("study_events").
			Unique(),
	}
}
