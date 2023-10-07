package wikitext

import (
	"fmt"
	"strings"
)

func ParseDefinition(raw string) string {
	tokens := tokenizer(raw)
	definitionParts := []string{}

	for _, t := range tokens {
		fmt.Println(t)
		switch t.Type {
		case "text":

		case "link":
			break
		}
	}

	definition := formatText(strings.Join(definitionParts, " "))

	// take first sentence as definition
	definition = strings.SplitAfter(definition, ".")[0]

	// TODO - Handle redirects more gracefully instead of removing outright
	if strings.Contains(definition, "#REDIRECT") || strings.Contains(definition, "#redirect") {
		return ""
	}

	return definition
}

// handles the different link types
func resolveLink(linkParts []string) []string {
	l := strings.Join(linkParts, " ")

	// category, interwiki link, or file
	if strings.Contains(l, ":") {
		return []string{}
	}

	// link with display text [[name of page|display text]]
	if strings.Contains(l, "|") {
		parts := strings.Split(l, "|")
		hasDisplay := len(parts) == 2

		if hasDisplay {
			return parts[1:]
		}
		return []string{}
	}

	return strings.Split(l, " ")
}

func formatText(def string) string {
	return strings.ReplaceAll(def, "  ", " ")
}
