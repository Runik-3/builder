package main

import (
	"flag"

	d "github.com/runik-3/builder/pkg/dict"
	l "github.com/runik-3/builder/pkg/lexicon"
	// output
	// definition depth (how many lines)
	// print
)

func main() {
	wikiUrl := flag.String("w", "", "wikiUrl")
	entryLimit := flag.Int("l", 10000, "limit")
	flag.Parse()

	BuildDictionary(wikiUrl, entryLimit)
}

func BuildDictionary(wikiUrl *string, entryLimit *int) {
	dict := d.Dict{Lex: l.New()}
	dict.GenerateDefinitionsFromWiki(wikiUrl, entryLimit)
	dict.Lex.Print()
}
