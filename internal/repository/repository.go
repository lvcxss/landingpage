package repository

import "landingpage/internal/data"

type PostRepository interface {
	GetByID(id string) (*data.Post, error)
	Insert(post *data.Post) error
	GetExperiences() ([]data.Experience, error)
	GetAll() ([]data.Post, error)
}
