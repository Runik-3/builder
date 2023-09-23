package builder

import (
	"flag"

	"github.com/runik-3/builder/pkg/dict"
	"github.com/runik-3/builder/pkg/wikiBot"
)

type Lexicon struct {
	Name string
	Dict *dict.Dict
}

func BuildDictionary() {
	wikiUrl := flag.String("w", "", "wikiUrl")
	entryLimit := flag.Int("l", 10000, "limit")
	flag.Parse()

	lex := Lexicon{Dict: dict.New()}
	lex.buildWords(wikiUrl, entryLimit)
	lex.Dict.Print()
	println(lex.Name)
}

func (d Lexicon) buildWords(wikiUrl *string, pageLimit *int) Lexicon {
	wikibot.GenerateDefinitionsFromWiki(d.Dict, wikiUrl, pageLimit)
	return d
}
