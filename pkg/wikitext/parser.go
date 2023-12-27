package wikitext

import (
	"strings"
)

func ParseDefinition(raw string, depth int) string {
	tokens := tokenizer(raw)
	definition := ""

	for _, t := range tokens {
		switch t.Type {
		case "text":
			definition += t.Value

		case "link":
			definition += resolveLink(t.Value)
		}
	}

	// resolve depth of definition
	sentences := strings.SplitAfter(definition, ".")[0:depth]
	definition = strings.Join(sentences, "")

	// remove doubled spaces
	definition = strings.ReplaceAll(definition, "  ", " ")
	// trim whitespace
	definition = strings.TrimSpace(definition)

	// TODO - Handle redirects more gracefully instead of removing outright
	if strings.Contains(definition, "#REDIRECT") || strings.Contains(definition, "#redirect") {
		return ""
	}

	return definition
}

// handles the different link types
func resolveLink(link string) string {
	// category, interwiki link, or file
	if strings.Contains(link, ":") {
		return ""
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
