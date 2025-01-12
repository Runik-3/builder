package wikitext

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Token struct {
	Type string
	//SubType string
	Value string
}

// Since we only care about text content, tokenized results are
// represented as a flat collection.
type TokenCollection []Token

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

type Tokenizer struct {
	tokens     TokenCollection
	state      State
	characters []string
}

func NewTokenizer(raw string) Tokenizer {
	cleaned := cleanDocument(raw)
	characters := strings.Split(cleaned, "")

	return Tokenizer{characters: characters, state: State{"text_start"}}
}

// returns a flat collection of tokens, but respects nested
// token types like templates.
func (t *Tokenizer) Tokenize() Tokenizer {
	i := 0
	for i < len(t.characters) {
		char := t.characters[i]
		tt := t.GetTokenType(&i)

		switch t.state.getState() {
		case "text_start":
			if tt == text {
				t.state = t.state.set("text")
				t.newToken(char, t.state)
				break
			}
			fallthrough
		case "text":
			if tt == link_start {
				t.state = t.state.set("link")
				t.newToken("", t.state)
				break
			}
			if tt == template_start {
				t.state = t.state.set("template")
				t.newToken("", t.state)
				break
			}
			if tt == table_start {
				t.state = t.state.set("table")
				t.newToken("", t.state)
				break
			}
			if tt == heading {
				t.state = t.state.set("heading")
				t.newToken("", t.state)
				break
			}
			if tt == text {
				t.appendToToken(char)
			}

		case "link":
			if tt == link_end {
				t.state = t.state.set("text_start")
				break
			}
			if tt == text {
				t.appendToToken(char)
				break
			}

		case "heading":
			if tt == heading {
				t.state = t.state.set("text_start")
				break
			}
			if tt == text {
				t.appendToToken(char)
				break
			}

		case "template":
			if tt == template_start {
				t.state = append(t.state, "template")
				break
			}
			if tt == template_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set("text_start")
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt == text {
				t.appendToToken(char)
				break
			}

		case "table":
			if tt == table_start {
				t.state = append(t.state, "table")
				break
			}
			if tt == table_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set("text_start")
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt == text {
				t.appendToToken(char)
				break
			}
		}
		i++
	}
	return *t

}

// this should return tokengrammar
func (t *Tokenizer) GetTokenType(i *int) TokenType {
	chars := t.characters
	// TODO: the tokenizer is in need of a cleanup, this can be improved.
	tknType := text

	currChar := chars[*i]
	// end of file
	if len(chars) == *i+1 {
		// FIXME: this is gross -- can this be solved with an EOF token?
		if currChar == "=" {
			tknType = heading
		}
		if currChar == "}" {
			tknType = table_end
		}
		return tknType
	}

	nextChar := chars[*i+1]

	if currChar == "[" && nextChar == "[" {
		*i++
		return link_start
	}
	if currChar == "]" && nextChar == "]" {
		*i++
		return link_end
	}
	if currChar == "{" && nextChar == "{" {
		*i++
		return template_start
	}
	if currChar == "}" && nextChar == "}" {
		*i++
		return template_end
	}
	if currChar == "{" {
		return table_start
	}
	if currChar == "}" {
		return table_end
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

func (t *Tokenizer) newToken(char string, state State) {
	newTkn := Token{Type: state.getState(), Value: char}
	t.tokens = append(t.tokens, newTkn)
}

func (t *Tokenizer) appendToToken(char string) {
	if len(t.tokens) > 0 {
		currTkn := &t.tokens[len(t.tokens)-1]
		// potential improvement would be to store this as an array until we finish
		// the token and then we can close it off by joining.
		currTkn.Value += char
	}
}

// Prepares document for tokenization.
func cleanDocument(t string) string {
	//s := cleanHtmlContent(t)
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

// FIXME: Regex isn't quite the right tool for this job. It's not grainular
// enough to handle nested or side-by-side tags. HTML tags will have to be
// handled within the tokenizer.

// removes html tags whose inner text we don't want to preserve
func cleanHtmlContent(s string) string {
	tags := []string{"ref"}

	for _, t := range tags {
		reg := regexp.MustCompile(fmt.Sprintf("<%s.*>.*<\\/%s>", t, t))
		s = reg.ReplaceAllString(s, "")
	}

	return s
}

// removes html tags from text, but preserves inner text
func cleanHtmlTags(s string) string {
	reg := regexp.MustCompile("<[^<>]*>")

	return reg.ReplaceAllString(s, "")
}
