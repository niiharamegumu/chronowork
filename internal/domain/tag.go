package domain

import "time"

// Tag represents a label for categorizing work entries.
type Tag struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
