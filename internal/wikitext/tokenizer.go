package wikitext

import (
	"encoding/json"
	"regexp"
	"strings"
)

type TokenType int

const (
	link_start TokenType = iota
	link_end
	table_start
	table_end
	heading
	text
)

// Since we only care about text content, tokenized results are
// represented as a flat collection.
type TokenCollection []Token

type Token struct {
	Type  string
	Value string
}

func (t *TokenCollection) Stringify() (string, error) {
	json, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

type State []string

func (s State) getState() string {
	return s[len(s)-1]
}

func (s State) pop() State {
	if len(s) > 0 {
		return s[0 : len(s)-1]
	}
	return s
}

func (s State) set(state string) State {
	return append(s.pop(), state)
}

func (s State) len() int {
	return len(s)
}

// returns a flat collection of tokens, but respects nested
// token types like tables.
func tokenizer(raw string) TokenCollection {
	cleaned := cleanDocument(raw)
	tokens := TokenCollection{}
	state := State{"text_start"}

	characters := strings.Split(cleaned, "")

	i := 0
	for i < len(characters) {
		char := characters[i]
		tt := tokenType(characters, &i)

		switch state.getState() {
		case "text_start":
			if tt == text {
				state = state.set("text")
				newToken(char, state, &tokens)
				break
			}
			fallthrough
		case "text":
			if tt == link_start {
				state = state.set("link")
				newToken("", state, &tokens)
				break
			}
			if tt == table_start {
				state = state.set("table")
				newToken("", state, &tokens)
				break
			}
			if tt == heading {
				state = state.set("heading")
				newToken("", state, &tokens)
				break
			}
			if tt == text {
				appendToToken(char, tokens)
			}

		case "link":
			if tt == link_end {
				state = state.set("text_start")
				break
			}
			if tt == text {
				appendToToken(char, tokens)
				break
			}

		case "heading":
			if tt == heading {
				state = state.set("text_start")
				break
			}
			if tt == text {
				appendToToken(char, tokens)
				break
			}

		case "table":
			if tt == table_start {
				state = append(state, "table")
				break
			}
			if tt == table_end {
				// check stack and assign appropriate state
				if state.len() == 1 {
					state = state.set("text_start")
				} else {
					appendToToken("", tokens) // close table
					state = state.pop()
				}
				break
			}
			if tt == text {
				appendToToken(char, tokens)
				break
			}
		}

		i++
	}
	return tokens
}

func tokenType(chars []string, i *int) TokenType {
	// TODO: the tokenizer is in need of a cleanup, this can be improved.
	tknType := text

	currChar := chars[*i]
	// end of file
	if len(chars) == *i+1 {
		if currChar == "=" {
			tknType = heading
		}
		return tknType
	}

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
	if currChar == "=" {
		for range chars[*i:] {
			if len(chars) > *i+1 && chars[*i+1] == "=" {
				*i++
				continue
			}
			tknType = heading
			break
		}
	}

	return tknType
}

func newToken(char string, state State, tokens *TokenCollection) {
	newTkn := Token{Type: state.getState(), Value: char}
	*tokens = append(*tokens, newTkn)
}
func appendToToken(char string, tokens TokenCollection) {
	if len(tokens) > 0 {
		currTkn := &tokens[len(tokens)-1]
		currTkn.Value += char
	}
}

// Prepares document for tokenization.
func cleanDocument(t string) string {
	s := cleanHtmlTags(t)

	// strip urls from text (likely came from <ref> tags that got cleaned above
	reg := regexp.MustCompile(`(f|ht)(tp)(s?)(://)(\S*)[.|/]([^\s\]\}]*)`)
	s = reg.ReplaceAllString(s, "")

	// wikitext bold
	s = strings.ReplaceAll(s, "'''", "")
	// wikitext italics
	s = strings.ReplaceAll(s, "''", "")

	return s
}

// removes html tags from text, but preserves inner text
func cleanHtmlTags(t string) string {
	reg := regexp.MustCompile("<[^<>]*>")
	s := reg.ReplaceAllString(t, "")

	return s
}
