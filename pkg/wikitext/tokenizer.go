package wikitext

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	link_start TokenType = iota
	link_end
	table_start
	table_end
	text
)

type Token struct {
	Type  string
	Value string
}

func tokenizer(raw string) []Token {
	cleaned := cleanText(raw)
	tokens := []Token{}
	state := "text_start"

	for _, t := range strings.Fields(cleaned) {
		tt := tokenType(t)

		switch tt {
		case text:
			if state != "text_start" && len(tokens) != 0 {
				appendToToken(tokens, cleanHtml(t))
			} else {
				state = "text"
				tokens = newToken(tokens, state, cleanHtml(t))
			}

		case link_start:
			if strings.Contains(state, "text") {
				state = "link"
				tokens = newToken(tokens, state, trimLinks(t))
				// links in tables should be appended as-is
			} else {
				appendToToken(tokens, t)
			}

		case link_end:
			if state == "link" {
				appendToToken(tokens, trimLinks(t))
				state = "text_start"
				// handle if link is one token eg. [[hi]]
			} else if strings.Contains(state, "text") {
				state = "link"
				tokens = newToken(tokens, state, trimLinks(t))
				state = "text_start"
			} else {
				appendToToken(tokens, t)
			}

		case table_start:
			if strings.Contains(state, "text") {
				state = "table"
				tokens = newToken(tokens, state, trimLinks(t))
				// tables inside tokens append as is.
			} else {
				appendToToken(tokens, t)
			}

		case table_end:
			if state == "table" {
				appendToToken(tokens, trimLinks(t))
				state = "text_start"
			} else if strings.Contains(state, "text") {
				state = "table"
				tokens = newToken(tokens, state, trimLinks(t))
				state = "text_start"
			} else {
				appendToToken(tokens, t)
			}
		}
	}

	return tokens
}

func newToken(tokens []Token, state string, content string) []Token {
	newTkn := Token{Type: state, Value: fmt.Sprintf("%s ", content)}
	return append(tokens, newTkn)
}
func appendToToken(tokens []Token, content string) {
	currTkn := &tokens[len(tokens)-1]
	currTkn.Value = currTkn.Value + fmt.Sprintf("%s ", content)
}

func tokenType(t string) TokenType {
	tknType := text

	if strings.Contains(t, "[[") {
		tknType = link_start
	}
	if strings.Contains(t, "]]") {
		tknType = link_end
	}
	if strings.Contains(t, "{{") {
		tknType = table_start
	}
	if strings.Contains(t, "}}") {
		tknType = table_end
	}
	return tknType
}

func trimTables(t string) string {
	s := strings.ReplaceAll(t, "{{", "")
	s = strings.ReplaceAll(s, "}}", "")
	return s
}

func trimLinks(t string) string {
	s := strings.ReplaceAll(t, "[[", "")
	s = strings.ReplaceAll(s, "]]", "")
	return s
}

func cleanText(t string) string {
	s := cleanHtml(t)
	// wikitext bold
	s = strings.ReplaceAll(s, "'''", "")
	// wikitext italics
	return strings.ReplaceAll(s, "''", "")
}

func cleanHtml(t string) string {
	isTag := false
	parts := strings.Split(t, "")
	cleaned := []string{}

	for _, p := range parts {
		if p == "<" {
			isTag = true
		}

		if !isTag {
			cleaned = append(cleaned, p)
		}

		if p == ">" {
			isTag = false
		}
	}

	return strings.Join(cleaned, "")
}
