package wikitext

import (
	"errors"
	"strings"
)

func ParseDefinition(raw string, depth int) (string, error) {
	if raw == "" {
		return "", errors.New("No page content.")
	}
	tokenizer := NewTokenizer(raw)
	definition := ""

	// While loop, while def length is less than sentence depth, keep batching
	batch := 1
	batchSize := 300
	for !isDefinitionParsed(&definition, &tokenizer, depth) {
		tokenizer.Tokenize(TokenizerOptions{batch, batchSize})
		definition = ""

		// FIXME: we're doing a little extra work by resetting the definition above
		// and processing tokens over again. Perhaps an unapply last token func?
		for _, t := range tokenizer.tokens {
			switch t.Type {
			case "text":
				definition += t.Value

			case "link":
				definition += resolveLink(t.Value)
			}
		}
		batch++
	}

	// resolve depth of definition
	sentences := strings.SplitAfter(definition, ". ")
	if depth <= len(sentences) {
		sentences = sentences[0:depth]
	}
	definition = strings.Join(sentences, "")

	// remove doubled spaces
	definition = strings.ReplaceAll(definition, "  ", " ")
	// trim whitespace
	definition = strings.TrimSpace(definition)

	// TODO - Handle redirects more gracefully instead of removing outright
	if strings.Contains(strings.ToLower(definition), "#redirect") {
		return "", nil
	}

	return definition, nil
}

func isDefinitionParsed(def *string, t *Tokenizer, depth int) bool {
	// If we have no tokens, we just started.
	if len(t.tokens) == 0 {
		return false
	}

	// resolve depth of definition
	sentences := strings.SplitAfter(*def, ". ")
	if depth < len(sentences) || t.tokens[len(t.tokens)-1].Type == "EOF" {
		return true
	}

	return false
}

// handles the different link types
func resolveLink(link string) string {
	// Since we're only parsing text from wikis, we don't want to handle link
	// types that don't include display text (ie. files, media, categories)
	unsupported := []string{"file:", "media:", "category:"}

	for _, u := range unsupported {
		if strings.Contains(strings.ToLower(link), u) {
			return ""
		}
	}

	// link with display text [[name of page|display text]]
	if strings.Contains(link, "|") {
		parts := strings.Split(link, "|")
		hasDisplay := len(parts) == 2

		if hasDisplay {
			return strings.Join(parts[1:], "")
		}
		return ""
	}

	return link
}
