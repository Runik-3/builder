package dict

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	l "github.com/runik-3/builder/pkg/lexicon"
	wikibot "github.com/runik-3/builder/pkg/wikiBot"
	"github.com/runik-3/builder/pkg/wikitext"
)

type Dict struct {
	Name string
	Lex  *l.Lexicon
}

func (d Dict) GenerateDefinitionsFromWiki(wikiUrl string, entryLimit int) {
	w := wikibot.CreateClient(wikiUrl)

	entries := 0

	// initial call has empty apfrom
	res := wikibot.GetWikiPageBatch(w, "", entryLimit)

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

		if entries == entryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}

		res = wikibot.GetWikiPageBatch(w, res.Continue.Apcontinue, entryLimit-entries)
	}

	fmt.Printf("📖 Found %d entries \n", entries)
}

// TODO - support for different formats
func (d Dict) Write(path string) string {
	p := path

	fileName := fmt.Sprintf("%s.json", d.Name)
	normalizedPath := filepath.Join(filepath.FromSlash(p), fileName)

	f, fileErr := os.Create(normalizedPath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	defer f.Close()

	json, marshalErr := json.Marshal(d.Lex)
	if fileErr != nil {
		log.Fatal(marshalErr)
	}

	_, writeErr := f.WriteString(string(json))
	if fileErr != nil {
		log.Fatal(writeErr)
	}

	fmt.Printf("Successfully built dictionary at %s\n", normalizedPath)

	return normalizedPath
}
