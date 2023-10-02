package wikitext

import (
	"strings"
)

func ParseDefinition(raw string) string {
	tokens := tokenizer(raw)
	definition := ""

	for _, t := range tokens {
		if t.Type == "text" {
			definition += t.Value
		}
		if t.Type == "link" {
			definition += resolveLink(t.Value)
		}
	}

	// take first sentence as definition
	definition = strings.SplitAfter(definition, ".")[0]

	return definition
}

// handles the different link types
func resolveLink(l string) string {
	// link with display text [[name of page|display text]]
	if strings.Contains(l, "|") {
		parts := strings.Split(l, "|")
		hasDisplay := len(parts) == 2

		if hasDisplay {
			return parts[1]
		}
		return ""
	}

	// category or interwiki link
	if strings.Contains(l, ":") {
		return ""
	}

	return l
}
