package wikitext

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
)

type Token struct {
	Type  State
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

type StateStack []State

func (s StateStack) getState() State {
	return s[len(s)-1]
}

func (s StateStack) pop() StateStack {
	if len(s) > 0 {
		return s[0 : len(s)-1]
	}
	return s
}

func (s StateStack) set(state State) StateStack {
	return append(s.pop(), state)
}

func (s StateStack) len() int {
	return len(s)
}

type Tokenizer struct {
	tokens     TokenCollection
	state      StateStack
	characters []string
	i          int
}

type TokenizerOptions struct {
	// If this value exists, the tokenizer will process the text in batches
	batchSize int
}

func NewTokenizer(raw string) Tokenizer {
	cleaned := cleanDocument(raw)
	characters := strings.Split(cleaned, "")

	return Tokenizer{characters: characters, state: StateStack{unset}}
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
	batch := len(t.characters)
	if options.batchSize != 0 {
		maybeBatch := t.i + options.batchSize
		if maybeBatch < batch {
			batch = maybeBatch
		}
	}

	for t.i < batch {
		char := t.characters[t.i]
		tt := t.GetTokenType()

		switch t.state.getState() {
		case unset:
			if tt.Token == text_token {
				t.state = t.state.set(text)
				t.newToken(char)
				break
			}
			fallthrough
		case text:
			if tt.Token != text_token {
				t.state = t.state.set(tt.State)
				t.newToken("")
				break
			}
			if tt.Token == text_token {
				t.appendToToken(char)
			}

		case link:
			if tt.Token == link_end {
				t.state = t.state.set(unset)
				break
			}
			if tt.Token == text_token {
				t.appendToToken(char)
				break
			}

		case heading:
			if tt.Token == heading_token {
				t.state = t.state.set(unset)
				break
			}
			if tt.Token == text_token {
				t.appendToToken(char)
				break
			}

		case template:
			if tt.Token == template_start {
				t.state = append(t.state, template)
				break
			}
			if tt.Token == template_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set(unset)
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt.Token == text_token {
				t.appendToToken(char)
				break
			}

		case table:
			if tt.Token == table_start {
				t.state = append(t.state, table)
				break
			}
			if tt.Token == table_end {
				// check stack and assign appropriate state
				if t.state.len() == 1 {
					t.state = t.state.set(unset)
				} else {
					t.state = t.state.pop()
				}
				break
			}
			if tt.Token == text_token {
				t.appendToToken(char)
				break
			}
		}

		// End of file
		if t.i == len(t.characters)-1 {
			t.state.set(EOF)
			t.newToken("")
		}

		t.i++
	}

	return *t
}

func (t *Tokenizer) GetTokenType() TokenGrammar {
	currChar := t.characters[t.i]

	matcherFunc, ok := matcherFunctions[currChar]
	if !ok {
		return TEXT_TOKEN
	}

	r, isMatch := matcherFunc(&t.i, &t.characters)
	if isMatch {
		return r
	}

	return TEXT_TOKEN
}

func (t *Tokenizer) newToken(char string) {
	newTkn := Token{Type: t.state.getState(), Value: char}
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
	s := cleanHtml(t)
	// strip urls from text (likely came from <ref> tags that got cleaned above
	reg := regexp.MustCompile(`(f|ht)(tp)(s?)(://)(\S*)[.|/]([^\s\]\}]*)`)
	s = reg.ReplaceAllString(s, "")

	// wikitext bold
	s = strings.ReplaceAll(s, "'''", "")
	// wikitext italics
	s = strings.ReplaceAll(s, "''", "")

	return s
}

func cleanHtml(s string) string {
	// Tags whose inner content we don't want to preserve
	removeHtmlContent := []string{"ref"}
	chars := strings.Split(s, "")
	cleanedText := []string{}

	tagType := ""
	// when true, write to cleanedText
	write := true
	tag := ""
	tags := []string{} // to track nesting of clean html content tags
	for i, c := range chars {
		if c == "<" {
			tagType = "open"

			lookAhead := canLookAhead(i, &chars)
			if lookAhead && chars[i+1] == "/" {
				tagType = "close"
			}
			write = false
			continue
		}

		if c == ">" {
			// ignore attributes that may have been captured along with the tag
			parsedTag := strings.Split(tag, " ")[0]
			// we don't want text contents from this tag type
			if tagType == "open" {
				write = true
				// If this tag is a remove html content tag don't write
				if slices.Contains(removeHtmlContent, parsedTag) {
					write = false
					tags = append(tags, tag)
				}
			} else if tagType == "close" {
				// There was an issue parsing the tag, the content may be malformed.
				// Let's still clean the malformed content. If this becomes a regular
				// pattern for pages with valid content, we can come up with a
				// strategy to handle this.
				if len(parsedTag) < 1 {
					tags = pop(tags)
				} else if slices.Contains(removeHtmlContent, parsedTag[1:]) {
					tags = pop(tags)
				}
				write = true
			}

			// If this tag is nested in a remove html content tag don't write
			if len(tags) > 0 {
				write = false
			}
			tagType = ""
			tag = ""
			continue
		}

		if write {
			cleanedText = append(cleanedText, c)
		}
		if tagType != "" {
			tag += c
		}
	}

	return strings.Join(cleanedText, "")
}

func pop(s []string) []string {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}
