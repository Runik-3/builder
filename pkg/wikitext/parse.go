package wikitext

import (
	"strings"
)

type TokenType int

const (
	link_start TokenType = iota
	link_end
	text
)

type Token struct {
	Type  string
	Value []string
}

func Parse(raw string) []Token {
	ts := tokenizer(raw)

	return ts
}

// consider a refactor where text is the first class type and everything else lives within text
// no need to get fancy with links, we just need start and end and state?
func tokenizer(raw string) []Token {
	tokens := []Token{}
	state := "text_start"

	for _, t := range strings.Split(raw, " ") {
		tt := tokenType(t)

		switch tt {
		case text:
			if state != "text_start" && len(tokens) != 0 {
				currTkn := &tokens[len(tokens)-1]
				currTkn.Value = append(currTkn.Value, t)
			} else {
				state = "text"
				newTkn := Token{Type: state, Value: []string{t}}
				tokens = append(tokens, newTkn)
			}

		case link_start:
			state = "link"
			newTkn := Token{Type: state, Value: []string{trimLinks(t)}}
			tokens = append(tokens, newTkn)

		case link_end:
			if state == "link" {
				currTkn := &tokens[len(tokens)-1]
				currTkn.Value = append(currTkn.Value, trimLinks(t))
			} else {
				state = "link"
				newTkn := Token{Type: state, Value: []string{trimLinks(t)}}
				tokens = append(tokens, newTkn)
			}

			state = "text_start"
		}
	}

	return tokens
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
