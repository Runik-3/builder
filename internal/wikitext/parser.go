package wikitext

import (
	"fmt"
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
	sentences := strings.SplitAfter(definition, ". ")
	sentences = NormalizeSentences(sentences)
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

// Enforces rules around what constitutes a sentence.
//   - A sentence ends on a `. `
//   - A sentence must contain at least 2 words (eliminates acronyms from
//     counting as a sentence, eg. a.k.a. should not be a sentence simply because
//     it includes a `. `)
func NormalizeSentences(sentences []string) []string {
	for i, sentence := range sentences {
		fmt.Println(sentence)
		words := strings.Split(sentence, " ")
		if len(words) <= 1 && i+1 < len(sentences) {
			sentences[i] = sentence + sentences[i+1]
			i++
			// remove next entry
			if i+2 < len(sentences) {
				sentences = append(sentences[:i+1], sentences[i+2:]...)
			}
		}
	}

	return sentences
}
