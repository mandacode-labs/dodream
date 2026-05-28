package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Notebook holds the schema definition for the Notebook entity.
type Notebook struct {
	ent.Schema
}

// Fields of the Notebook.
func (Notebook) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Notebook.
func (Notebook) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).
			Ref("notebooks").
			Unique().
			Required(),
		edge.To("notebook_cards", NotebookCard.Type),
	}
}
