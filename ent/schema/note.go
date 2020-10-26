package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// Note holds the schema definition for the Note entity.
type Note struct {
	ent.Schema
}

// Fields of the Note.
func (Note) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("content"),
	}
}

// Edges of the Note.
func (Note) Edges() []ent.Edge {
	return nil
}
