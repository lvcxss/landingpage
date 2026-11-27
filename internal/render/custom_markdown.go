package render

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// CustomParseMD transforms raw Markdown text into terminal-styled HTML.
// It demonstrates how a custom markdown renderer works from scratch,
// allowing complete control over emitted HTML tags, CSS classes, and terminal decorations.
func CustomParseMD(source string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	var out strings.Builder

	inCodeBlock := false
	codeBlockLang := ""
	var codeBlockLines []string

	inBlockquote := false
	var blockquoteLines []string

	inList := false
	isOrderedList := false

	inTable := false
	var tableHeader []string
	var tableRows [][]string

	flushBlockquote := func() {
		if inBlockquote && len(blockquoteLines) > 0 {
			content := strings.Join(blockquoteLines, "<br>\n")
			out.WriteString(fmt.Sprintf("<blockquote>\n%s\n</blockquote>\n", parseInline(content)))
			blockquoteLines = nil
			inBlockquote = false
		}
	}

	flushList := func() {
		if inList {
			if isOrderedList {
				out.WriteString("</ol>\n")
			} else {
				out.WriteString("</ul>\n")
			}
			inList = false
		}
	}

	flushTable := func() {
		if inTable && len(tableHeader) > 0 {
			out.WriteString("<table>\n<thead>\n<tr>\n")
			for _, th := range tableHeader {
				out.WriteString(fmt.Sprintf("  <th>%s</th>\n", parseInline(strings.TrimSpace(th))))
			}
			out.WriteString("</tr>\n</thead>\n<tbody>\n")
			for _, row := range tableRows {
				out.WriteString("<tr>\n")
				for _, cell := range row {
					out.WriteString(fmt.Sprintf("  <td>%s</td>\n", parseInline(strings.TrimSpace(cell))))
				}
				out.WriteString("</tr>\n")
			}
			out.WriteString("</tbody>\n</table>\n")
			tableHeader = nil
			tableRows = nil
			inTable = false
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 1. Fenced Code Blocks (```lang ... ```)
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				// Close code block with custom terminal window styling
				escapedCode := html.EscapeString(strings.Join(codeBlockLines, "\n"))
				titleLang := codeBlockLang
				if titleLang == "" {
					titleLang = "sh"
				}

				out.WriteString("<div class=\"terminal-window markdown-code-block\">\n")
				out.WriteString("  <div class=\"terminal-window-bar\">\n")
				out.WriteString("    <div class=\"terminal-window-dots\">\n")
				out.WriteString("      <span class=\"terminal-dot dot-red\"></span>\n")
				out.WriteString("      <span class=\"terminal-dot dot-yellow\"></span>\n")
				out.WriteString("      <span class=\"terminal-dot dot-green\"></span>\n")
				out.WriteString("    </div>\n")
				out.WriteString(fmt.Sprintf("    <span class=\"terminal-window-title\">$ cat script.%s</span>\n", html.EscapeString(titleLang)))
				out.WriteString("  </div>\n")
				out.WriteString(fmt.Sprintf("  <pre><code class=\"language-%s\">%s</code></pre>\n", html.EscapeString(titleLang), escapedCode))
				out.WriteString("</div>\n")

				inCodeBlock = false
				codeBlockLang = ""
				codeBlockLines = nil
				continue
			} else {
				flushBlockquote()
				flushList()
				flushTable()
				inCodeBlock = true
				codeBlockLang = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				codeBlockLines = nil
				continue
			}
		}

		if inCodeBlock {
			codeBlockLines = append(codeBlockLines, line)
			continue
		}

		// Empty lines reset block elements
		if trimmed == "" {
			flushBlockquote()
			flushList()
			flushTable()
			continue
		}

		// 2. Horizontal Rule (--- or ***)
		if (trimmed == "---" || trimmed == "***" || trimmed == "___") && len(trimmed) >= 3 {
			flushBlockquote()
			flushList()
			flushTable()
			out.WriteString("<hr>\n")
			continue
		}

		// 3. Headings (# H1, ## H2, etc.)
		if strings.HasPrefix(trimmed, "#") {
			flushBlockquote()
			flushList()
			flushTable()

			level := 0
			for level < len(trimmed) && trimmed[level] == '#' {
				level++
			}
			if level <= 6 && len(trimmed) > level && trimmed[level] == ' ' {
				headingText := strings.TrimSpace(trimmed[level:])
				parsedHeading := parseInline(headingText)

				if level == 1 {
					out.WriteString(fmt.Sprintf("<h1>%s</h1>\n", parsedHeading))
				} else if level == 2 {
					out.WriteString(fmt.Sprintf("<h2><span class=\"text-primary\">$</span> %s</h2>\n", parsedHeading))
				} else if level == 3 {
					out.WriteString(fmt.Sprintf("<h3><span class=\"text-primary\">#</span> %s</h3>\n", parsedHeading))
				} else {
					out.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, parsedHeading, level))
				}
				continue
			}
		}

		// 4. Blockquotes (> quote)
		if strings.HasPrefix(trimmed, ">") {
			flushList()
			flushTable()
			inBlockquote = true
			quoteContent := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
			blockquoteLines = append(blockquoteLines, quoteContent)
			continue
		} else if inBlockquote {
			flushBlockquote()
		}

		// 5. Tables (| col1 | col2 |)
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			flushBlockquote()
			flushList()

			rawCells := strings.Split(trimmed, "|")
			var cells []string
			for i := 1; i < len(rawCells)-1; i++ {
				cells = append(cells, rawCells[i])
			}

			// Check if this is separator line |---|---|
			isSeparator := true
			for _, cell := range cells {
				cleanCell := strings.Trim(strings.TrimSpace(cell), ":- ")
				if cleanCell != "" {
					isSeparator = false
					break
				}
			}

			if isSeparator {
				// Separator line, do not append to data
				continue
			}

			if !inTable {
				inTable = true
				tableHeader = cells
			} else {
				tableRows = append(tableRows, cells)
			}
			continue
		} else if inTable {
			flushTable()
		}

		// 6. Lists (Unordered - / * and Ordered 1. )
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushBlockquote()
			flushTable()
			if !inList || isOrderedList {
				flushList()
				inList = true
				isOrderedList = false
				out.WriteString("<ul>\n")
			}
			itemText := parseInline(strings.TrimSpace(trimmed[2:]))
			out.WriteString(fmt.Sprintf("  <li>%s</li>\n", itemText))
			continue
		}

		if isOrderedListPrefix(trimmed) {
			flushBlockquote()
			flushTable()
			if !inList || !isOrderedList {
				flushList()
				inList = true
				isOrderedList = true
				out.WriteString("<ol>\n")
			}
			dotIdx := strings.Index(trimmed, ".")
			itemText := parseInline(strings.TrimSpace(trimmed[dotIdx+1:]))
			out.WriteString(fmt.Sprintf("  <li>%s</li>\n", itemText))
			continue
		}

		if inList {
			flushList()
		}

		// 7. Regular Paragraph
		out.WriteString(fmt.Sprintf("<p>%s</p>\n", parseInline(trimmed)))
	}

	flushBlockquote()
	flushList()
	flushTable()

	return out.String(), nil
}

// parseInline processes inline elements: bold, italic, code, links, images
func parseInline(text string) string {
	// 1. Inline code `code` (placeholder to prevent further formatting inside code)
	codeRegex := regexp.MustCompile("`([^`]+)`")
	var codeTokens []string
	text = codeRegex.ReplaceAllStringFunc(text, func(m string) string {
		match := codeRegex.FindStringSubmatch(m)
		if len(match) > 1 {
			token := fmt.Sprintf("@@INLINECODETOKEN%d@@", len(codeTokens))
			codeTokens = append(codeTokens, fmt.Sprintf("<code>%s</code>", html.EscapeString(match[1])))
			return token
		}
		return m
	})

	// 2. Images ![alt](url)
	imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = imgRegex.ReplaceAllString(text, `<img src="$2" alt="$1">`)

	// 3. Links [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllString(text, `<a href="$2">$1</a>`)

	// 4. Bold **text** or __text__
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = boldRegex.ReplaceAllString(text, `<strong>$1</strong>`)
	boldUnderRegex := regexp.MustCompile(`__([^_]+)__`)
	text = boldUnderRegex.ReplaceAllString(text, `<strong>$1</strong>`)

	// 5. Italic *text* or _text_
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	text = italicRegex.ReplaceAllString(text, `<em>$1</em>`)
	italicUnderRegex := regexp.MustCompile(`_([^_]+)_`)
	text = italicUnderRegex.ReplaceAllString(text, `<em>$1</em>`)

	// 6. Restore inline code
	for i, tokenVal := range codeTokens {
		tokenKey := fmt.Sprintf("@@INLINECODETOKEN%d@@", i)
		text = strings.ReplaceAll(text, tokenKey, tokenVal)
	}

	return text
}

func isOrderedListPrefix(s string) bool {
	if len(s) < 3 {
		return false
	}
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	return i > 0 && i < len(s) && s[i] == '.' && (i+1 < len(s) && s[i+1] == ' ')
}
