// Code generated by entc, DO NOT EDIT.

package note

const (
	// Label holds the string label denoting the note type in the database.
	Label = "note"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"

	// Table holds the table name of the note in the database.
	Table = "notes"
)

// Columns holds all SQL columns for note fields.
var Columns = []string{
	FieldID,
	FieldContent,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
