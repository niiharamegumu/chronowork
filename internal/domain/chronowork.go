// Package domain contains pure domain entities without external dependencies.
package domain

import "time"

// ChronoWork represents a work tracking entry.
type ChronoWork struct {
	ID            uint
	Title         string
	ProjectTypeID uint
	TagID         uint
	StartTime     time.Time
	EndTime       time.Time
	IsTracking    bool
	TotalSeconds  int
	Confirmed     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships (loaded when needed)
	ProjectType *ProjectType
	Tag         *Tag
}
