package lexicon

import (
	"fmt"

	"github.com/runik-3/builder/pkg/dict"
	"github.com/runik-3/builder/pkg/wikiBot"
)

type Lexicon struct {
	Name string
	Dict *dict.Dict
}

func (l Lexicon) GenerateDefinitionsFromWiki(wikiUrl *string, entryLimit *int) {
	w := wikibot.CreateClient(*wikiUrl)

	entries := 0

	// initial call has empty apfrom
	res := wikibot.GetWikiPages(w, "", *entryLimit)

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			// parsing content happens here
			l.Dict.Add(dict.Entry{Word: p.Title, Definition: p.Revisions[0].Slots.Main.Content})
			entries++
		}

		if entries == *entryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}

		res = wikibot.GetWikiPages(w, res.Continue.Apcontinue, *entryLimit-entries)
	}

	fmt.Printf("📖 Found %d entries \n", entries)
}
