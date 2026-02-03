package model

// Bookmark represents a user saved link in the system.
// It stores a human readable description, the original URL,
// a short code that can be used for quick access, and the owner user.
// The struct is mapped to the "bookmarks" table in the database using GORM tags.
//
// Fields:
//   - ID: Inherited from Base, unique identifier (UUID) for the bookmark
//   - Description: Optional text describing the bookmark
//   - URL: Original URL to be stored and accessed
//   - Code: Short, unique code generated for the bookmark
//   - UserID: Foreign key referencing the owner user
//   - User: Preloaded user entity for relational queries
type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code"`
	UserID      string `json:"-" gorm:"type:uuid;column:user_id"`
	User        User   `gorm:"references:ID" json:"-"`
}
