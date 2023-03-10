// Code generated by ent, DO NOT EDIT.

package category

const (
	// Label holds the string label denoting the category type in the database.
	Label = "category"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldIcon holds the string denoting the icon field in the database.
	FieldIcon = "icon"
	// Table holds the table name of the category in the database.
	Table = "category"
)

// Columns holds all SQL columns for category fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldIcon,
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