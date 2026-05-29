package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// State holds the schema definition for the State entity.
type State struct {
	ent.Schema
}

// Fields of the State.
func (State) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.Time("next_review_at"),
		field.Float("interval"),
		field.Float("ease_factor"),
		field.Int("total_reviews").Default(0),
		field.Int("total_successful").Default(0),
		field.Int("streak").Default(0),
		field.Int("reviews_last_1d").Default(0),
		field.Int("reviews_last_3d").Default(0),
		field.Int("reviews_last_7d").Default(0),
		field.Float("avg_response_time_ms").Optional().Nillable(),
		field.Time("last_review_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the State.
func (State) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("states").
			Unique().
			Required(),
		edge.From("card", Card.Type).
			Ref("states").
			Unique().
			Required(),
	}
}
