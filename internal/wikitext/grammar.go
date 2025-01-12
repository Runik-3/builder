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
	html_start
	html_end
	text
	EOF
)

type TokenGrammar struct {
	Token TokenType
	State string
}

var grammar = map[string]TokenGrammar{
	"[[":     {Token: link_start, State: "link"},
	"]]":     {Token: link_end, State: "link"},
	"{{":     {Token: template_start, State: "template"},
	"}}":     {Token: template_end, State: "template"},
	"{":      {Token: table_start, State: "table"},
	"}":      {Token: table_end, State: "table"},
	"=":      {Token: heading, State: "heading"},
	"==":     {Token: heading, State: "heading"},
	"===":    {Token: heading, State: "heading"},
	"====":   {Token: heading, State: "heading"},
	"=====":  {Token: heading, State: "heading"},
	"======": {Token: heading, State: "heading"},
	"<%>":    {Token: html_start, State: "html"},
	"</%>":   {Token: html_end, State: "html"},
}
