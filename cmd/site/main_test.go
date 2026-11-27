package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublishIndependentContent(t *testing.T) {
	writeups, blog := t.TempDir(), t.TempDir()
	add := func(root, path, text string) {
		t.Helper()
		target := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	cvRepo := t.TempDir()
	add(cvRepo, "cv.yaml", `name: Test Portfolio
education: PDF-specific section
experiences:
  - name: Timeline entry
    start: "2026"
    description: Visible description
  - name: PDF-only entry
    show_on_site: false
`)
	add(writeups, "notes/demo.md", "---\ndate: 2025-01-02\n---\n# A writeup\n")
	add(blog, "notes/demo.md", "---\r\ntitle: Blog title\r\ndate: 2026-01-02\r\n---\r\nHello there.")
	add(writeups, "scripts/omit.md", "ignored")
	add(writeups, "notes/omit.py", "ignored")
	for _, base := range []string{"", "/landingpage"} {
		out := t.TempDir()
		if err := build(filepath.Join(cvRepo, "cv.yaml"), writeups, blog, out, base); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(out, "index.html"))
		if err != nil {
			t.Fatal(err)
		}
		html := string(b)
		for _, expected := range []string{
			`href="` + base + `/static/style.css"`, `href="` + base + `/post/notes-demo"`,
			`href="` + base + `/post/project-notes-demo"`, `data-category="project"`, "Blog title", "A writeup", "Test Portfolio", "Visible description",
		} {
			if !strings.Contains(html, expected) {
				t.Errorf("missing %q", expected)
			}
		}
		if strings.Contains(html, "PDF-only entry") || strings.Contains(html, "PDF-specific section") {
			t.Error("CV-only content in site")
		}
		if strings.Contains(html, "omit") {
			t.Error("unexpected article")
		}
		for _, rel := range []string{"post/notes-demo/index.html", "post/project-notes-demo/index.html", "static/terminus.css"} {
			if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
				t.Error(err)
			}
		}
	}
}
