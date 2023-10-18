package dict

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	f "github.com/runik-3/builder/pkg/formatter"
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
func (d Dict) Write(path string, format string) string {
	p := path

	fileName := fmt.Sprintf("%s.%s", d.Name, format)
	normalizedPath := filepath.Join(filepath.FromSlash(p), fileName)

	file, fileErr := os.Create(normalizedPath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	defer file.Close()

	formatted := f.Format(format, *d.Lex)

	_, writeErr := file.WriteString(formatted)
	if fileErr != nil {
		log.Fatal(writeErr)
	}

	fmt.Printf("Successfully built dictionary at %s\n", normalizedPath)

	return normalizedPath
}
