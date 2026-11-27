package render

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithXHTML()),
	)
}

// ParseMD parses markdown into terminal-themed HTML.
// It uses CustomParseMD for terminal styling (terminal window bars, dots, $ prefixes),
// and falls back to Goldmark if needed.
func ParseMD(source string) (string, error) {
	result, err := CustomParseMD(source)
	if err == nil && result != "" {
		return result, nil
	}
	return ParseMDGoldmark(source)
}

// ParseMDGoldmark converts markdown using Goldmark GFM.
func ParseMDGoldmark(source string) (string, error) {
	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
