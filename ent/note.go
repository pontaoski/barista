// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/appadeia/barista/ent/note"
	"github.com/facebook/ent/dialect/sql"
)

// Note is the model entity for the Note schema.
type Note struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Content holds the value of the "content" field.
	Content string `json:"content,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Note) scanValues() []interface{} {
	return []interface{}{
		&sql.NullString{}, // id
		&sql.NullString{}, // content
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Note fields.
func (n *Note) assignValues(values ...interface{}) error {
	if m, n := len(values), len(note.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field id", values[0])
	} else if value.Valid {
		n.ID = value.String
	}
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field content", values[0])
	} else if value.Valid {
		n.Content = value.String
	}
	return nil
}

// Update returns a builder for updating this Note.
// Note that, you need to call Note.Unwrap() before calling this method, if this Note
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *Note) Update() *NoteUpdateOne {
	return (&NoteClient{config: n.config}).UpdateOne(n)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (n *Note) Unwrap() *Note {
	tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: Note is not a transactional entity")
	}
	n.config.driver = tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *Note) String() string {
	var builder strings.Builder
	builder.WriteString("Note(")
	builder.WriteString(fmt.Sprintf("id=%v", n.ID))
	builder.WriteString(", content=")
	builder.WriteString(n.Content)
	builder.WriteByte(')')
	return builder.String()
}

// Notes is a parsable slice of Note.
type Notes []*Note

func (n Notes) config(cfg config) {
	for _i := range n {
		n[_i].config = cfg
	}
}
