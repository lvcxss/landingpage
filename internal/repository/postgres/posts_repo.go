package postgres

import (
	"database/sql"
	"errors"
	"time"

	"landingpage/internal/data"
)

type PostRepo struct {
	DB *sql.DB
}

func (repo *PostRepo) GetByID(id string) (*data.Post, error) {
	query := `SELECT id, title, file_path, COALESCE(image_url, ''), post_type, created_at FROM posts WHERE id = $1 OR file_path LIKE '%' || $1 OR title = $1 LIMIT 1`
	var post data.Post
	var pTypeStr string
	err := repo.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.FilePath, &post.ImageURL, &pTypeStr, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	post.PostType, _ = data.ParsePostType(pTypeStr)
	return &post, nil
}

func (repo *PostRepo) Insert(post *data.Post) error {
	postTypeStr := post.PostType.String()
	post.CreatedAt = time.Now()
	if post.ImageURL == "" {
		query := `INSERT INTO posts (title, file_path, post_type, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
		return repo.DB.QueryRow(query, post.Title, post.FilePath, postTypeStr, post.CreatedAt).Scan(&post.ID)
	}

	query := `INSERT INTO posts (title, file_path, image_url, post_type, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return repo.DB.QueryRow(query, post.Title, post.FilePath, post.ImageURL, postTypeStr, post.CreatedAt).Scan(&post.ID)
}

func (repo *PostRepo) GetExperiences() ([]data.Experience, error) {
	query := `SELECT id, start_date, end_date, "name", description, entity, url FROM experiences ORDER BY start_date DESC`
	var experiences []data.Experience

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var exp data.Experience
		err = rows.Scan(&exp.ID, &exp.StartDate, &exp.EndDate, &exp.Name, &exp.Description, &exp.Company, &exp.URL)
		if err != nil {
			return nil, err
		}
		experiences = append(experiences, exp)
	}
	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("couldn't find any experiences")
		}
		return nil, err
	}
	return experiences, nil
}

func (repo *PostRepo) GetAll() ([]data.Post, error) {
	query := `SELECT id, title, file_path, post_type, created_at FROM posts ORDER BY created_at DESC`

	var posts []data.Post

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post data.Post
		var pTypeStr string

		err = rows.Scan(&post.ID, &post.Title, &post.FilePath, &pTypeStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.PostType, _ = data.ParsePostType(pTypeStr)

		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
