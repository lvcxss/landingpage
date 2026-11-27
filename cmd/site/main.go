// Command site renders the portfolio to an empty directory for static hosting.
package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"landingpage/content"
	"landingpage/internal/data"
	"landingpage/internal/render"
	"landingpage/internal/ui"
)

type page struct {
	BasePath    string
	Title       string
	Content     template.HTML
	Experiences []data.Experience
	Timeline    []content.Experience
	Posts       []data.Post
}

func main() {
	cvPath := flag.String("cv", content.DefaultCVPath, "YAML file from the CV repository")
	writeups := flag.String("writeups", "docs/writeups", "writeups repository")
	blog := flag.String("blog", "", "optional blog repository")
	output := flag.String("out", "public", "new or empty output directory")
	base := flag.String("base", "", "URL prefix, e.g. /landingpage")
	flag.Parse()
	if err := build(*cvPath, *writeups, *blog, *output, strings.TrimRight(*base, "/")); err != nil {
		log.Fatal(err)
	}
}

func build(cvPath, writeups, blog, output, base string) error {
	entries, err := os.ReadDir(output)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if len(entries) != 0 {
		return fmt.Errorf("output directory %s must be empty", output)
	}
	cv, err := content.ReadCV(cvPath)
	if err != nil {
		return err
	}
	renderer, err := render.NewHTMLRenderer(ui.HTMLFiles)
	if err != nil {
		return err
	}
	home := page{BasePath: base, Title: cv.Name, Timeline: cv.Experiences}
	save := func(rel, tmpl string, p page) error {
		path := filepath.Join(output, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		return errors.Join(renderer.Execute(f, tmpl, p), f.Close())
	}
	sources := []struct {
		root     string
		category data.PostTypes
	}{{writeups, data.WriteUp}, {blog, data.Project}}
	seen := make(map[string]string)
	for _, source := range sources {
		if source.root == "" {
			continue
		}
		err := fs.WalkDir(os.DirFS(source.root), ".", func(rel string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				if rel != "." && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "scripts") {
					return fs.SkipDir
				}
				return nil
			}
			if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(rel), ".md") {
				return nil
			}
			post, body, err := loadPost(source.root, rel, source.category)
			if err != nil {
				return fmt.Errorf("%s: %w", rel, err)
			}
			if previous, exists := seen[post.ID]; exists {
				return fmt.Errorf("duplicate URL %s for %s and %s", post.ID, previous, rel)
			}
			seen[post.ID] = rel
			html, err := render.ParseMD(body)
			if err != nil {
				return err
			}
			home.Posts = append(home.Posts, post)
			return save(filepath.Join("post", post.ID, "index.html"), "post.html", page{Title: post.Title, BasePath: base, Content: template.HTML(html)})
		})
		if err != nil {
			return err
		}
	}
	sort.SliceStable(home.Posts, func(i, j int) bool { return home.Posts[i].CreatedAt.After(home.Posts[j].CreatedAt) })
	if err := save("index.html", "index.html", home); err != nil {
		return err
	}
	return os.CopyFS(filepath.Join(output, "static"), ui.StaticFiles)
}

var separator = regexp.MustCompile(`[^\p{L}\p{N}]+`)

func loadPost(root, rel string, category data.PostTypes) (data.Post, string, error) {
	var post data.Post
	source, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return post, "", err
	}
	body := strings.ReplaceAll(string(source), "\r\n", "\n")
	var meta struct {
		Title string
		Date  string
	}
	if rest, ok := strings.CutPrefix(body, "---\n"); ok {
		header, article, ok := strings.Cut(rest, "\n---\n")
		if !ok {
			return post, "", fmt.Errorf("unclosed YAML header")
		}
		if err := yaml.Unmarshal([]byte(header), &meta); err != nil {
			return post, "", err
		}
		body = article
	}
	stem := strings.TrimSuffix(rel, filepath.Ext(rel))
	slug := strings.Trim(separator.ReplaceAllString(strings.ToLower(stem), "-"), "-")
	if slug == "" {
		return post, "", fmt.Errorf("filename needs a letter or digit")
	}
	if category == data.Project {
		slug = "project-" + slug
	}
	if meta.Title == "" {
		meta.Title = filepath.Base(stem)
		for _, line := range strings.Split(body, "\n") {
			if title, ok := strings.CutPrefix(line, "# "); ok {
				meta.Title = strings.TrimSpace(title)
				break
			}
		}
	}
	var date time.Time
	if meta.Date != "" {
		date, err = time.Parse("2006-01-02", meta.Date)
	} else {
		// Use commit history so a fresh checkout does not change article dates.
		var stamp []byte
		stamp, err = exec.Command("git", "-C", root, "log", "-1", "--format=%aI", "--", rel).Output()
		if err == nil {
			date, err = time.Parse(time.RFC3339, strings.TrimSpace(string(stamp)))
		}
		if err != nil {
			return post, "", fmt.Errorf("set date in the YAML header or commit the article for a Git date: %w", err)
		}
	}
	return data.Post{ID: slug, Title: meta.Title, CreatedAt: date, PostType: category}, body, err
}
