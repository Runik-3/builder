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

type TokenizerOptions struct {
	// If these values exist, the tokenizer will process the text in batches
	batch     int
	batchSize int
}

func NewTokenizer(raw string) Tokenizer {
	cleaned := cleanDocument(raw)
	characters := strings.Split(cleaned, "")

	return Tokenizer{characters: characters, state: State{"text_start"}}
}

func (t *Tokenizer) batcher(batch int, size int) []string {
	start := (batch - 1) * size
	end := batch * size

	if start > len(t.characters) {
		return []string{}
	}

	if end > len(t.characters) {
		end = len(t.characters)
	}

	return t.characters[start:end]
}

// FIXME: this can be cleaned up, let's make it leaner and easier to extend --
// already started the improvements with the tokentype func, let's keep going
func (t *Tokenizer) Tokenize(options TokenizerOptions) Tokenizer {
	chars := t.characters
	i := 0
	batchStart := 0

	if options.batch != 0 && options.batchSize != 0 {
		chars = t.batcher(options.batch, options.batchSize)
		batchStart = (options.batch - 1) * options.batchSize
		i = batchStart
	}

	for i < batchStart+len(chars) {
		char := t.characters[i]
		tt := t.GetTokenType(&i)

		switch t.state.getState() {
		case "text_start":
			if tt.Token == text {
				t.state = t.state.set("text")
				t.newToken(char, t.state)
				break
			}
			fallthrough
		case "text":
			if tt.Token == link_start {
				t.state = t.state.set("link")
				t.newToken("", t.state)
				break
			}
			if tt.Token == template_start {
				t.state = t.state.set("template")
				t.newToken("", t.state)
				break
			}
			if tt.Token == table_start {
				t.state = t.state.set("table")
				t.newToken("", t.state)
				break
			}
			if tt.Token == heading {
				t.state = t.state.set("heading")
				t.newToken("", t.state)
				break
			}
			if tt.Token == text {
				t.appendToToken(char)
			}

		case "link":
			if tt.Token == link_end {
				t.state = t.state.set("text_start")
				break
			}
			if tt.Token == text {
				t.appendToToken(char)
				break
			}

		case "heading":
			if tt.Token == heading {
				t.state = t.state.set("text_start")
				break
			}
			if tt.Token == text {
				t.appendToToken(char)
				break
			}

		case "template":
			if tt.Token == template_start {
				t.state = append(t.state, "template")
				break
			}
			if tt.Token == template_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set("text_start")
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt.Token == text {
				t.appendToToken(char)
				break
			}

		case "table":
			if tt.Token == table_start {
				t.state = append(t.state, "table")
				break
			}
			if tt.Token == table_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set("text_start")
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt.Token == text {
				t.appendToToken(char)
				break
			}
		}
		i++
	}

	return *t
}

func (t *Tokenizer) GetTokenType(i *int) TokenGrammar {
	currChar := t.characters[*i]

	matcherFunc, ok := matcherFunctions[currChar]
	if !ok {
		return TokenGrammar{Token: text, State: "text"}
	}

	r, isMatch := matcherFunc(i, &t.characters)
	if isMatch {
		return r
	}

	return TokenGrammar{Token: text, State: "text"}
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
// handled within the tokenizer. -- THIS IS DISABLED FOR NOW AND NEEDS TO BE
// RE-ENABLED BEFORE NEXT RELEASE
//
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
