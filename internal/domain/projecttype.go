package domain

import "time"

// ProjectType represents a project category with associated tags.
type ProjectType struct {
	ID        uint
	Name      string
	Tags      []Tag
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetTagNames returns the names of all associated tags.
func (p *ProjectType) GetTagNames() []string {
	var tagNames []string
	for _, tag := range p.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}
