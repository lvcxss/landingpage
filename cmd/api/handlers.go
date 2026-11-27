package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"landingpage/content"
	"landingpage/internal/data"
	"landingpage/internal/render"
)

type PageData struct {
	BasePath    string
	Timeline    []content.Experience
	Title       string
	RawTitle    string
	Content     template.HTML
	Experiences []data.Experience
	Posts       []data.Post
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	experiences, err := app.posts.GetExperiences()
	if err != nil {
		if err.Error() == "couldn't find any experiences" {
			experiences = []data.Experience{}
		}
	}

	posts, err := app.posts.GetAll()
	if err != nil {
		app.logger.Println("Aviso ao buscar posts da BD:", err)
		posts = []data.Post{}
	}

	// Merge posts found directly on disk in docs/posts
	posts = mergePostsFromDisk(posts, app.logger)

	cv, err := content.ReadCV(content.DefaultCVPath)
	if os.IsNotExist(err) {
		cv.Name = "Lucas Oliveira"
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{
		Timeline:    cv.Experiences,
		Title:       cv.Name,
		Experiences: experiences,
		Posts:       posts,
	}

	app.renderer.Render(w, http.StatusOK, "index.html", data)
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		id = r.PathValue("name")
	}
	if id == "" {
		http.NotFound(w, r)
		return
	}

	var post *data.Post
	var err error

	// 1. Attempt database lookup
	if app.posts != nil {
		post, err = app.posts.GetByID(id)
	}

	// 2. If not found in DB or DB error, attempt disk lookup in docs/posts/
	if post == nil || err != nil {
		post = findPostOnDisk(id)
		if post == nil {
			app.logger.Printf("Post não encontrado nem na BD nem no disco para ID: %s", id)
			http.NotFound(w, r)
			return
		}
	}

	// 3. Read markdown file
	filePath := post.FilePath
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fallbackPath := filepath.Join(".", "docs", "posts", filepath.Base(filePath))
		fileBytes, err = os.ReadFile(fallbackPath)
		if err != nil {
			app.logger.Printf("Erro a ler ficheiro do post %s: %v", filePath, err)
			http.NotFound(w, r)
			return
		}
	}

	// 4. Render markdown using our terminal-styled parser
	parsed, err := render.ParseMD(string(fileBytes))
	if err != nil {
		app.logger.Printf("Erro ao processar markdown para %s: %v", post.Title, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:    post.Title,
		RawTitle: post.Title,
		Content:  template.HTML(parsed),
	}
	app.renderer.Render(w, http.StatusOK, "post.html", data)
}

func (app *application) syncPost(w http.ResponseWriter, r *http.Request) {
	totpCode := r.Header.Get("X-TOTP-Code")
	secretBase32 := app.config.AdminToken
	valid, err := totp.ValidateCustom(totpCode, secretBase32, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil || !valid {
		app.logger.Println("Acesso negado: Código TOTP inválido")
		http.Error(w, "Unauthorized - Invalid TOTP", http.StatusUnauthorized)
		return
	}

	var input struct {
		Title    string `json:"title"`
		FilePath string `json:"file_path"`
		PostType string `json:"post_type"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fileName := filepath.Base(input.FilePath)
	if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md"
	}
	serverPath := filepath.Join(".", "docs", "posts", fileName)

	if err := os.MkdirAll(filepath.Dir(serverPath), 0o755); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(serverPath, []byte(input.Content), 0o644); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pType, _ := data.ParsePostType(input.PostType)
	post := &data.Post{
		Title:    input.Title,
		FilePath: serverPath,
		PostType: pType,
	}

	if err := app.posts.Insert(post); err != nil {
		app.logger.Printf("Aviso: Falha ao inserir post na BD (%v), mas guardado com sucesso em disco", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post sincronizado de forma segura!"))
}

func getPostsDir() string {
	candidates := []string{
		filepath.Join(".", "docs", "posts"),
		filepath.Join("..", "docs", "posts"),
		filepath.Join("..", "..", "docs", "posts"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return filepath.Join(".", "docs", "posts")
}

// Helper: merge posts found in docs/posts/ with any existing posts
func mergePostsFromDisk(existing []data.Post, logger *log.Logger) []data.Post {
	known := make(map[string]bool)
	for _, p := range existing {
		known[p.ID] = true
		known[filepath.Base(p.FilePath)] = true
		known[strings.TrimSuffix(filepath.Base(p.FilePath), ".md")] = true
	}

	postsDir := getPostsDir()
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		return existing
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		baseID := strings.TrimSuffix(entry.Name(), ".md")
		if known[baseID] || known[entry.Name()] {
			continue
		}

		filePath := filepath.Join(postsDir, entry.Name())
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		content := string(contentBytes)
		title := extractTitleFromMD(content, strings.Title(strings.ReplaceAll(baseID, "-", " ")))
		pType := detectPostType(entry.Name(), content)

		info, _ := entry.Info()
		modTime := time.Now()
		if info != nil {
			modTime = info.ModTime()
		}

		post := data.Post{
			ID:        baseID,
			Title:     title,
			FilePath:  filePath,
			PostType:  pType,
			CreatedAt: modTime,
		}

		existing = append(existing, post)
		known[baseID] = true
	}

	return existing
}

// Helper: locate a markdown post file on disk by slug or filename
func findPostOnDisk(id string) *data.Post {
	cleanID := filepath.Clean(id)
	cleanID = filepath.Base(cleanID) // Prevent directory traversal

	postsDir := getPostsDir()
	candidates := []string{
		filepath.Join(postsDir, cleanID+".md"),
		filepath.Join(postsDir, cleanID),
	}

	for _, candidate := range candidates {
		contentBytes, err := os.ReadFile(candidate)
		if err == nil {
			content := string(contentBytes)
			baseID := strings.TrimSuffix(filepath.Base(candidate), ".md")
			title := extractTitleFromMD(content, strings.Title(strings.ReplaceAll(baseID, "-", " ")))
			pType := detectPostType(candidate, content)
			info, _ := os.Stat(candidate)
			modTime := time.Now()
			if info != nil {
				modTime = info.ModTime()
			}
			return &data.Post{
				ID:        baseID,
				Title:     title,
				FilePath:  candidate,
				PostType:  pType,
				CreatedAt: modTime,
			}
		}
	}

	// Scan directory for partial match
	entries, err := os.ReadDir(postsDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			baseID := strings.TrimSuffix(entry.Name(), ".md")
			if strings.EqualFold(baseID, cleanID) {
				filePath := filepath.Join(postsDir, entry.Name())
				contentBytes, err := os.ReadFile(filePath)
				if err == nil {
					content := string(contentBytes)
					title := extractTitleFromMD(content, strings.Title(strings.ReplaceAll(baseID, "-", " ")))
					pType := detectPostType(entry.Name(), content)
					info, _ := entry.Info()
					modTime := time.Now()
					if info != nil {
						modTime = info.ModTime()
					}
					return &data.Post{
						ID:        baseID,
						Title:     title,
						FilePath:  filePath,
						PostType:  pType,
						CreatedAt: modTime,
					}
				}
			}
		}
	}

	return nil
}

func extractTitleFromMD(content string, defaultTitle string) string {
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return defaultTitle
}

func detectPostType(filename, content string) data.PostTypes {
	lower := strings.ToLower(filename + " " + content)
	if strings.Contains(lower, "write-up") || strings.Contains(lower, "writeup") || strings.Contains(lower, "sqli") || strings.Contains(lower, "ctf") || strings.Contains(lower, "exploit") {
		return data.WriteUp
	}
	if strings.Contains(lower, "tutorial") || strings.Contains(lower, "guide") {
		return data.Guide
	}
	if strings.Contains(lower, "project") {
		return data.Project
	}
	return data.WriteUp
}
