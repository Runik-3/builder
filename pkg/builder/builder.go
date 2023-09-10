package builder

import (
	"flag"

	"github.com/runik-3/builder/pkg/dict"
	"github.com/runik-3/builder/pkg/wikiBot"
)

type Lexicon struct {
	Dict *dict.Dict
}

func BuildDictionary() {
	wikiUrl := flag.String("w", "", "wikiUrl")
	entryLimit := flag.Int("l", 100000, "limit")
	flag.Parse()

	lex := Lexicon{Dict: dict.New()}
	lex.words(wikiUrl, entryLimit)
	lex.Dict.Print()
}

func (d Lexicon) words(wikiUrl *string, pageLimit *int) Lexicon {
	wikibot.GenerateWordList(d.Dict, wikiUrl, pageLimit)
	return d
}
