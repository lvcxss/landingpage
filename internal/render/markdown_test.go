package render

import (
	"os"
	"strings"
	"testing"
)

func TestParseMD(t *testing.T) {
	sample := "# Hello World\n\nThis is **bold** and `code`.\n\n```python\nprint('hello')\n```"
	res, err := CustomParseMD(sample)
	if err != nil {
		t.Fatalf("CustomParseMD failed: %v", err)
	}

	if !strings.Contains(res, "<h1>Hello World</h1>") {
		t.Errorf("Expected h1, got: %s", res)
	}
	if !strings.Contains(res, "<strong>bold</strong>") {
		t.Errorf("Expected bold, got: %s", res)
	}
	if !strings.Contains(res, "<code>code</code>") {
		t.Errorf("Expected code, got: %s", res)
	}
	if !strings.Contains(res, "terminal-window") {
		t.Errorf("Expected terminal-window code block, got: %s", res)
	}
}

func TestPrioritisePost(t *testing.T) {
	bytes, err := os.ReadFile("../../docs/posts/prioritise.md")
	if err != nil {
		t.Skip("docs/posts/prioritise.md not found, skipping")
	}
	res, err := CustomParseMD(string(bytes))
	if err != nil {
		t.Fatalf("Failed to parse prioritise.md: %v", err)
	}
	if !strings.Contains(res, "Prioritise") {
		t.Errorf("Expected Prioritise in output")
	}
}
