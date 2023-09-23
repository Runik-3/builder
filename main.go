package main

import (
	"flag"

	d "github.com/runik-3/builder/pkg/dict"
	l "github.com/runik-3/builder/pkg/lexicon"
)

func main() {
	wikiUrl := flag.String("w", "", "wikiUrl")
	entryLimit := flag.Int("l", 10000, "limit")
	flag.Parse()

	BuildDictionary(wikiUrl, entryLimit)
}

func BuildDictionary(wikiUrl *string, entryLimit *int) {
	lex := l.Lexicon{Dict: d.New()}
	lex.GenerateDefinitionsFromWiki(wikiUrl, entryLimit)
	lex.Dict.Print()
}
