package domain

import "time"

// Setting represents application configuration.
type Setting struct {
	ID                 uint
	RelativeDate       uint
	PersonDay          uint
	DisplayAsPersonDay bool
	DownloadPath       string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
