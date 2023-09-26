package dict

import (
	"fmt"

	l "github.com/runik-3/builder/pkg/lexicon"
	"github.com/runik-3/builder/pkg/wikiBot"
	"github.com/runik-3/builder/pkg/wikitext"
)

type Dict struct {
	Name string
	Lex  *l.Lexicon
}

func (d Dict) GenerateDefinitionsFromWiki(wikiUrl *string, entryLimit *int) {
	w := wikibot.CreateClient(*wikiUrl)

	entries := 0

	// initial call has empty apfrom
	res := wikibot.GetWikiPages(w, "", *entryLimit)

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			// parsing content happens here
			fDef := wikitext.Parse(p.Revisions[0].Slots.Main.Content)

			d.Lex.Add(l.Entry{Word: p.Title, Definition: fDef})
			entries++
		}

		if entries == *entryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}

		// call this get batch or something?
		res = wikibot.GetWikiPages(w, res.Continue.Apcontinue, *entryLimit-entries)
	}

	fmt.Printf("📖 Found %d entries \n", entries)
}
