package builder

import (
	"flag"

	"github.com/runik-3/builder/pkg/dict"
	"github.com/runik-3/builder/pkg/wikiBot"
)

type Lexicon struct {
	RawDict *dict.Dict
}

func BuildDictionary() {
	wikiUrl := flag.String("u", "", "wikiUrl")
	// low default limit of 5 for testing, should be 500
	pageLimit := flag.Int("p", 5, "pageLimit")
	flag.Parse()

	lex := Lexicon{RawDict: dict.New()}
	lex.words(wikiUrl, pageLimit)
}

func (d Lexicon) words(wikiUrl *string, pageLimit *int) Lexicon {
	wikibot.GenerateWordList(d.RawDict, wikiUrl, pageLimit)
	return d
}
