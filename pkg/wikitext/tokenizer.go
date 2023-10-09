package wikitext

import (
	"encoding/json"
	"log"
	"regexp"
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

type TokenCollection []Token

type Token struct {
	Type  string
	Value string
}

func (t TokenCollection) Stringify() string {
	json, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func tokenizer(raw string) TokenCollection {
	cleaned := cleanDocument(raw)
	tokens := TokenCollection{}
	state := "text_start"

	characters := strings.Split(cleaned, "")

	i := 0
	for i < len(characters) {
		char := characters[i]
		tt := tokenType(characters, &i)

		switch state {
		case "text_start":
			if tt == text {
				state = "text"
				newToken(char, state, &tokens)
				break
			}
			fallthrough
		case "text":
			if tt == link_start {
				state = "link"
				newToken("", state, &tokens)
				break
			}
			if tt == table_start {
				state = "table"
				newToken("", state, &tokens)
				break
			}
			if tt == text {
				appendToToken(char, state, tokens)
			}

		case "link":
			if tt == link_end {
				state = "text_start"
				break
			}
			if tt == text {
				appendToToken(char, state, tokens)
				break
			}

		case "table":
			if tt == table_end {
				state = "text_start"
				break
			}
			if tt == text {
				appendToToken(char, state, tokens)
				break
			}
		}

		i++
	}
	return tokens
}

func tokenType(chars []string, i *int) TokenType {
	tknType := text

	// end of file
	if len(chars) == *i+1 {
		return tknType
	}

	currChar := chars[*i]
	nextChar := chars[*i+1]

	if currChar == "[" && nextChar == "[" {
		*i++
		tknType = link_start
	}
	if currChar == "]" && nextChar == "]" {
		*i++
		tknType = link_end
	}
	if currChar == "{" && nextChar == "{" {
		*i++
		tknType = table_start
	}
	if currChar == "}" && nextChar == "}" {
		*i++
		tknType = table_end
	}

	return tknType
}

func newToken(char string, state string, tokens *TokenCollection) {
	newTkn := Token{Type: state, Value: char}
	*tokens = append(*tokens, newTkn)
}
func appendToToken(char string, state string, tokens TokenCollection) {
	if len(tokens) > 0 {
		currTkn := &tokens[len(tokens)-1]
		currTkn.Value += char
	}
}

// Prepares document for tokenization.
func cleanDocument(t string) string {
	s := cleanHtml(t)

	// strip urls from text (likely came from <ref> tags that got cleaned above
	reg := regexp.MustCompile(`(f|ht)(tp)(s?)(://)(\S*)[.|/]([^\s\]\}]*)`)
	s = reg.ReplaceAllString(s, "")

	// wikitext bold
	s = strings.ReplaceAll(s, "'''", "")
	// wikitext italics
	return strings.ReplaceAll(s, "''", "")
}

// removes html tags from text
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
