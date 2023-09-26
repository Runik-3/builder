package wikitext

import (
	"strings"
)

type Token struct {
	Type  string
	Value []string
}

func Parse(raw string) []Token {
	ts := tokenizer(raw)

	return ts
}

func tokenizer(raw string) []Token {
	tokens := []Token{}
	tkn := Token{} // no need for tkn here? can just append to last token in list?
	// Constant write stream, change state based on start/end tokens?

	for _, t := range strings.Split(raw, " ") {
		tt, tag := tokenType(t)

		// consider a refactor where text is the first class type and everything else lives within text
		// no need to get fancy with links, we just need start and end and state?
		switch tt {
		case "link":
			if tag == "start" || tag == "" {
				// add prev tkn and clear tkn value
				tokens = append(tokens, tkn)
				tkn.Value = []string{}
			}

			tkn.Type = tt
			tkn.Value = append(tkn.Value, trimLinks(t))

			if tag == "end" || tag == "" {
				tokens = append(tokens, tkn)
				tkn.Value = []string{}
			}

		case "text":
			if tkn.Type != "text" {
				tkn.Type = tt
			}
			tkn.Value = append(tkn.Value, t)
		}
	}

	return tokens
}

func tokenType(t string) (string, string) {
	tknType := "text" // text, link, table
	tag := ""         // start, end, or empty to denote both

	if strings.Contains(t, "[[") {
		tknType = "link"
		tag = "start"
	}
	if strings.Contains(t, "]]") {
		tknType = "link"
		if tag == "start" {
			tag = ""
		} else {
			tag = "end"
		}
	}
	return tknType, tag
}

func trimTables(t string) string {
	s := strings.ReplaceAll(t, "{{", "")
	s = strings.ReplaceAll(s, "}}", "")
	return s
}

func trimLinks(t string) string {
	s := strings.ReplaceAll(t, "[[", "")
	s = strings.ReplaceAll(s, "]]", "")
	return s
}

func trimBold(line string) string {
	return strings.ReplaceAll(line, "'''", "")
}

func trimItalic(line string) string {
	return strings.ReplaceAll(line, "''", "")
}
