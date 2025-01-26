package wikitext

import (
	"strings"
)

func ParseDefinition(raw string, depth int) string {
	tokenizer := NewTokenizer(raw)
	tokenizer.Tokenize(TokenizerOptions{})

	definition := ""

	for _, t := range tokenizer.tokens {
		switch t.Type {
		case "text":
			definition += t.Value

		case "link":
			definition += resolveLink(t.Value)
		}
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
		return ""
	}

	return definition
}

// handles the different link types
func resolveLink(link string) string {
	// Since we're only parsing text from wikis, we don't want to handle link
	// types that don't include display text (ie. files, media, categories)
	unsupported := []string{"file:", "media:", "categrory:"}

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
