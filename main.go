package main

import (
	"flag"

	"github.com/runik-3/builder/pkg/builder"
	"github.com/runik-3/builder/pkg/dict"
)

func main() {
	BuildDictionary()
}

func BuildDictionary() {
	wikiUrl := flag.String("w", "", "wikiUrl")
	entryLimit := flag.Int("l", 10000, "limit")
	flag.Parse()

	lex := builder.Lexicon{Dict: dict.New()}
	lex.GenerateDefinitionsFromWiki(wikiUrl, entryLimit)
	lex.Dict.Print()
	println(lex.Name)
}
