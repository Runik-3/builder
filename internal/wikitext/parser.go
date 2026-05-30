package wikitext

import (
	"bytes"
	"errors"
	"html"
	"regexp"
	"slices"
	"strings"
)

func CleanWord(raw string) string {
	// Replace `"` char with `'` in headword
	return strings.ReplaceAll(raw, "\"", "'")
}

func ParseDefinition(raw string, depth int) (string, error) {
	tokenizer := NewTokenizer(raw)
	if len(tokenizer.characters) == 0 {
		return "", errors.New("No page content.")
	}
	byteDef := []byte{}

	// While loop, while def length is less than sentence depth, keep batching
	batchSize := 300
	for !isDefinitionParsed(&byteDef, &tokenizer, depth) {
		tokenizer.Tokenize(TokenizerOptions{batchSize})
		byteDef = []byte{}

		// FIXME: we're doing a little extra work by resetting the definition above
		// and processing tokens over again. Perhaps an unapply last token func?
		for _, t := range tokenizer.tokens {
			switch t.Type {
			case "text":
				byteDef = append(byteDef, t.Value...)
			case "link":
				byteDef = append(byteDef, resolveLink(t.Value)...)
			}
		}
	}

	// resolve depth of definition
	sentences := splitSentence(byteDef)
	if depth <= len(sentences) {
		sentences = sentences[0:depth]
	}
	definition := string(bytes.Join(sentences, []byte{}))

	// TODO: consider using a strings.Replacer for these replace ops. The post
	// and pre tokenization cleaning steps should each only iterate once through
	// the definition.

	// remove doubled spaces
	definition = strings.ReplaceAll(definition, "  ", " ")
	// trim whitespace
	definition = strings.TrimSpace(definition)
	definition = collapseNewlines(definition)
	definition = handleIndents(definition)
	// convert html escaped chars to ascii
	definition = html.UnescapeString(definition)

	return definition, nil
}

func isDefinitionParsed(def *[]byte, t *Tokenizer, depth int) bool {
	// If we have no tokens, we just started.
	if len(t.tokens) == 0 {
		return false
	}

	sentences := splitSentence(*def)
	// resolve depth of definition
	if depth < len(sentences) || t.tokens[len(t.tokens)-1].Type == "EOF" {
		return true
	}

	return false
}

// Handles the different link types
//
// If link is in our list of unsupported, ignore it
// If link type is specified and does not contain display, ignore it
// If link type is specified but does contain display, keep it
func resolveLink(link []byte) []byte {
	// Since we're only parsing text from wikis, we don't want to handle link
	// types that don't include useful display text (ie. files, media)
	unsupported := [][]byte{[]byte("file:"), []byte("media:")}

	for _, u := range unsupported {
		if bytes.Contains(bytes.ToLower(link), u) {
			return []byte{}
		}
	}

	// link with display text (pipe link) -- [[name of page|display text]]
	if slices.Contains(link, '|') {
		idx := slices.Index(link, '|')
		hasDisplayText := len(link) - 1 > idx
		if hasDisplayText {	
			return link[idx + 1:]
		}

		return []byte{}
	}

	if slices.Contains(link, ':') {
		return []byte{}
	}

	return link
}

var newlineReg = regexp.MustCompile(`(\n){3,}`)

// Collapses excessive newlines into double-spaced break.
func collapseNewlines(s string) string {
	return newlineReg.ReplaceAllString(s, "\n\n")
}

var indentReg = regexp.MustCompile(`(^|\n)(:){1,6}`)

// TODO indents should be handled by the tokenizer
func handleIndents(s string) string {
	// Matches up to 6 levels of indent on first line or any new line
	return indentReg.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ReplaceAll(match, ":", "  ")
	})
}

// Note: this is very overfit for the english language
func splitSentence(raw []byte) [][]byte {
	split := [][]byte{}
	splitIdx := 0

	for i, char := range raw {
		// EOF
		if !canLookAhead(i, &raw) {
			split = append(split, raw[splitIdx:])
			continue	
		}

		if char != '.' && char != '!' && char != '?' {
			continue	
		}

		peek := raw[i + 1]
		if peek != ' ' && peek != '\n' {
			continue
		}
 
		delimiterOffset := 2
		split = append(split, raw[splitIdx:i+delimiterOffset])
		splitIdx = i + delimiterOffset
	}

	return split
}
