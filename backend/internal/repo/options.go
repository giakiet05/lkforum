package repo

// Filter is a map to define query conditions.
// It's an alias for map[string]interface{} for better readability in function signatures.
// Example: Filter{"community_id": id, "is_deleted": false}
type Filter map[string]interface{}

// UpdateDocument is a map that defines an update operation, often using MongoDB operators.
// It's an alias for map[string]interface{} for better readability.
// Example: UpdateDocument{"$set": {"title": "New Title"}}
type UpdateDocument map[string]interface{}

// FindOptions defines generic options for find operations, independent of the database driver.
type FindOptions struct {
	Sort  map[string]int // Example: {"created_at": -1, "votes": 1}
	Skip  int64
	Limit int64
}
