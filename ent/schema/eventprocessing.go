package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// EventProcessing holds the schema definition for the EventProcessing entity.
type EventProcessing struct {
	ent.Schema
}

// Fields of the EventProcessing.
func (EventProcessing) Fields() []ent.Field {
	return []ent.Field{
		field.String("event_id").Unique().Immutable(),
		field.String("status").Default("pending"),
		field.String("processor_id").Optional().Nillable(),
		field.Time("started_at").Optional().Nillable(),
		field.Time("completed_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the EventProcessing.
func (EventProcessing) Edges() []ent.Edge {
	return nil
}
