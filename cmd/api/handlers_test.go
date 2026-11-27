package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"landingpage/internal/data"
	"landingpage/internal/render"
	"landingpage/internal/ui"
)

// mockRepo satisfies repository.PostRepository without requiring live Postgres
type mockRepo struct {
	posts []data.Post
}

func (m *mockRepo) GetByID(id string) (*data.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) Insert(post *data.Post) error {
	m.posts = append(m.posts, *post)
	return nil
}

func (m *mockRepo) GetExperiences() ([]data.Experience, error) {
	return []data.Experience{}, nil
}

func (m *mockRepo) GetAll() ([]data.Post, error) {
	return m.posts, nil
}

func setupTestApp(t *testing.T) *application {
	t.Helper()
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("docs/posts", 0755); err != nil {
		t.Fatal(err)
	}
	article := "# Prioritise\n\nBlind SQL Injection\n\n```text\nsample\n```\n"
	if err := os.WriteFile("docs/posts/prioritise.md", []byte(article), 0644); err != nil {
		t.Fatal(err)
	}
	renderer, err := render.NewHTMLRenderer(ui.HTMLFiles, "*.html")
	if err != nil {
		t.Fatalf("Failed to initialize HTMLRenderer: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	return &application{
		logger:   logger,
		renderer: renderer,
		posts:    &mockRepo{},
	}
}

func TestHomeDiscoversPostsFromDisk(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	app.home(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "prioritise") && !strings.Contains(body, "Prioritise") {
		t.Errorf("Expected homepage to list 'Prioritise' post from docs/posts, got body:\n%s", body)
	}
}

func TestShowPostRendersMarkdown(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest("GET", "/post/prioritise", nil)
	req.SetPathValue("id", "prioritise")
	rr := httptest.NewRecorder()

	app.showPost(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Prioritise") {
		t.Errorf("Expected rendered markdown to contain 'Prioritise', got:\n%s", body)
	}
	if !strings.Contains(body, "Blind SQL Injection") {
		t.Errorf("Expected rendered markdown to contain 'Blind SQL Injection', got:\n%s", body)
	}
	if !strings.Contains(body, "terminal-window") && !strings.Contains(body, "<pre>") {
		t.Errorf("Expected rendered markdown to contain styled code block, got:\n%s", body)
	}
}
