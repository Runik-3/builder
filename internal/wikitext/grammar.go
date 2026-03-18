package wikitext

type TokenType string

const (
	link_start     TokenType = "link_start"
	link_end       TokenType = "link_end"
	template_start TokenType = "template_start"
	template_end   TokenType = "template_end"
	table_start    TokenType = "table_start"
	table_end      TokenType = "table_end"
	// parse all headings (h1-h6) as a single rule since we don't preserve them
	heading_token TokenType = "heading_token"
	text_token    TokenType = "text_token"
)

type State string

const (
	link     State = "link"
	template State = "template"
	table    State = "table"
	heading  State = "heading"
	unset    State = "unset"
	text     State = "text"
	EOF      State = "EOF"
)

type TokenGrammar struct {
	Token TokenType
	State State
}

var TEXT_TOKEN TokenGrammar = TokenGrammar{Token: text_token, State: text}

var grammar = map[string]TokenGrammar{
	"[[":     {Token: link_start, State: link},
	"]]":     {Token: link_end, State: link},
	"{{":     {Token: template_start, State: template},
	"}}":     {Token: template_end, State: template},
	"{":      {Token: table_start, State: table},
	"}":      {Token: table_end, State: table},
	"======": {Token: heading_token, State: heading},
	"=====":  {Token: heading_token, State: heading},
	"====":   {Token: heading_token, State: heading},
	"===":    {Token: heading_token, State: heading},
	"==":     {Token: heading_token, State: heading},
	"=":      {Token: heading_token, State: heading},
}

var matcherFunctions = map[byte]func(*int, *[]byte) (TokenGrammar, bool){
	byte('['): matchNext,
	byte(']'): matchNext,
	byte('{'): matchThisOrNext,
	byte('}'): matchThisOrNext,
	byte('='): matchMany,
}

// Matches double characters in sequence (ie. `[[`)
func matchNext(i *int, chars *[]byte) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return TokenGrammar{}, false
	}
	next := (*chars)[*i+1]

	if curr != next {
		return TokenGrammar{}, false
	}

	return getTokenMatch([]byte{curr, next}, i)
}

func matchThisOrNext(i *int, chars *[]byte) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return getTokenMatch([]byte{curr}, i)
	}
	next := (*chars)[*i+1]
	if curr != next {
		return getTokenMatch([]byte{curr}, i)
	}

	return matchNext(i, chars)
}

func matchMany(i *int, chars *[]byte) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return getTokenMatch([]byte{curr}, i)
	}

	matches := []byte{curr}
	match := true
	j := 1
	for match {
		if !canLookAhead(*i+j, chars) {
			matches = append(matches, curr)
			match = false
			break
		}
		next := (*chars)[*i+j]
		if curr != next {
			match = false
			break
		}
		matches = append(matches, curr)
		j++
	}

	return getTokenMatch(matches, i)
}

func canLookAhead(i int, chars *[]byte) bool {
	return len(*chars) > i+1
}

func canLookBehind(i int) bool {
	return i-1 >= 0
}

func getTokenMatch(key []byte, i *int) (TokenGrammar, bool) {
	rule, ok := grammar[string(key)]
	if !ok {
		return TokenGrammar{}, false
	}

	*i += len(key) - 1
	return rule, true
}
