package wikitext

type TokenType int

const (
	link_start TokenType = iota
	link_end
	template_start
	template_end
	table_start
	table_end
	// parse all headings (h1-h6) as a single rule since we skip them anyway
	heading
	text
)

type TokenGrammar struct {
	Token TokenType
	State string
}

var TEXT_TOKEN TokenGrammar = TokenGrammar{Token: text, State: "text"}

var grammar = map[string]TokenGrammar{
	"[[":     {Token: link_start, State: "link"},
	"]]":     {Token: link_end, State: "link"},
	"{{":     {Token: template_start, State: "template"},
	"}}":     {Token: template_end, State: "template"},
	"{":      {Token: table_start, State: "table"},
	"}":      {Token: table_end, State: "table"},
	"======": {Token: heading, State: "heading"},
	"=====":  {Token: heading, State: "heading"},
	"====":   {Token: heading, State: "heading"},
	"===":    {Token: heading, State: "heading"},
	"==":     {Token: heading, State: "heading"},
	"=":      {Token: heading, State: "heading"},
}

var matcherFunctions = map[string]func(*int, *[]string) (TokenGrammar, bool){
	"[": matchNext,
	"]": matchNext,
	"{": matchThisOrNext,
	"}": matchThisOrNext,
	"=": matchMany,
}

// Matches double characters in sequence (ie. `[[`)
func matchNext(i *int, chars *[]string) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return TokenGrammar{}, false
	}
	next := (*chars)[*i+1]

	if curr != next {
		return TokenGrammar{}, false
	}

	return getTokenMatch(curr+next, i)
}

func matchThisOrNext(i *int, chars *[]string) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return getTokenMatch(curr, i)
	}
	next := (*chars)[*i+1]
	if curr != next {
		return getTokenMatch(curr, i)
	}

	return matchNext(i, chars)
}

func matchMany(i *int, chars *[]string) (TokenGrammar, bool) {
	curr := (*chars)[*i]
	if !canLookAhead(*i, chars) {
		return getTokenMatch(curr, i)
	}

	matches := curr
	match := true
	j := 1
	for match {
		if !canLookAhead(*i+j, chars) {
			matches += curr
			match = false
			break
		}
		next := (*chars)[*i+j]
		if curr != next {
			match = false
			break
		}
		matches += curr
		j++
	}

	return getTokenMatch(matches, i)
}

func canLookAhead(i int, chars *[]string) bool {
	return len(*chars) > i+1
}

func getTokenMatch(key string, i *int) (TokenGrammar, bool) {
	rule, ok := grammar[key]
	if !ok {
		return TokenGrammar{}, false
	}

	*i += len(key) - 1
	return rule, true
}
