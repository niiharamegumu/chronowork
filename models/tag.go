package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"size:255; required; unique" json:"name"`
}

func FindALlTags(db *gorm.DB) []Tag {
	var tags []Tag
	result := db.Find(&tags).Order("id desc")
	if result.Error != nil {
		return []Tag{}
	}
	return tags
}

func FindByTagId(db *gorm.DB, id uint) (Tag, error) {
	var tag Tag
	result := db.First(&tag, id)
	if result.Error != nil {
		return Tag{}, result.Error
	}
	return tag, nil
}

func AllTagNames(db *gorm.DB) []string {
	var tags []Tag
	result := db.Find(&tags)
	if result.Error != nil {
		return []string{}
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}

func TagsByNames(db *gorm.DB, names []string) []Tag {
	var tags []Tag
	result := db.Where("name IN ?", names).Find(&tags)
	if result.Error != nil {
		return []Tag{}
	}
	return tags
}
