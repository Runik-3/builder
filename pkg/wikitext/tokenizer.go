package wikitext

import (
	"strings"
)

type TokenType int

const (
	link_start TokenType = iota
	link_end
	table_start
	table_end
	html_start
	html_end
	text
)

type Token struct {
	Type  string
	Value []string
}

func tokenizer(raw string) []Token {
	tokens := []Token{}
	state := "text_start"

	for _, t := range strings.Fields(raw) {
		tt := tokenType(t)

		switch tt {
		case text:
			if state != "text_start" && len(tokens) != 0 {
				appendToToken(tokens, t)
			} else {
				state = "text"
				tokens = newToken(tokens, state, t)
			}

		case link_start:
			state = "link"
			tokens = newToken(tokens, state, trimLinks(t))

		case link_end:
			if state == "link" {
				appendToToken(tokens, trimLinks(t))
			} else {
				state = "link"
				tokens = newToken(tokens, state, trimLinks(t))
			}

			state = "text_start"
		}
	}

	return tokens
}

func newToken(tokens []Token, state string, content string) []Token {
	newTkn := Token{Type: state, Value: []string{content}}
	return append(tokens, newTkn)
}
func appendToToken(tokens []Token, content string) {
	currTkn := &tokens[len(tokens)-1]
	currTkn.Value = append(currTkn.Value, content)
}

func tokenType(t string) TokenType {
	tknType := text

	if strings.Contains(t, "[[") {
		tknType = link_start
	}
	if strings.Contains(t, "]]") {
		tknType = link_end
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

func trimBold(line string) string {
	return strings.ReplaceAll(line, "'''", "")
}

func trimItalic(line string) string {
	return strings.ReplaceAll(line, "''", "")
}
