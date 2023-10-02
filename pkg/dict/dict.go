package dict

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
			def := wikitext.ParseDefinition(p.Revisions[0].Slots.Main.Content)
			if def != "" {
				d.Lex.Add(l.Entry{Word: p.Title, Definition: def})
				entries++
			}
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

// TODO - support for different formats
func (d Dict) Write(path string) {
	p := path

	if p == "" {
		p = "dict.json"
	}

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	json, err := json.Marshal(d.Lex)
	f.WriteString(string(json))
}
