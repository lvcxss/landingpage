package data

import (
	"errors"
	"time"
)

type PostTypes int

const (
	ExperiencePost PostTypes = iota
	WriteUp
	Project
	Guide
)

var postStrings = map[PostTypes]string{
	ExperiencePost: "experience",
	WriteUp:        "write-up",
	Project:        "project",
	Guide:          "guide",
}

func (p PostTypes) String() string {
	val, ok := postStrings[p]
	if ok {
		return val
	}
	return "post"
}

func ParsePostType(s string) (PostTypes, error) {
	for key, val := range postStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, errors.New("invalid post type")
}

type Post struct {
	ID        string
	Title     string
	FilePath  string
	ImageURL  string
	PostType  PostTypes
	CreatedAt time.Time
}

type Experience struct {
	ID          string
	StartDate   time.Time
	EndDate     *time.Time
	Company     string
	Name        string
	Description string
	URL         *string
}
