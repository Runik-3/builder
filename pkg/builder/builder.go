package builder

import (
	"flag"

	wikibot "github.com/runik-3/builder/pkg/wikiBot"
)

func BuildDictionary() {
	wikiUrl := flag.String("u", "", "wikiUrl")
	// low default limit of 5 for testing, should be 500
	pageLimit := flag.Int("p", 5, "pageLimit")
	flag.Parse()

	wikibot.GenerateWordList(wikiUrl, pageLimit)
}
