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
	entries := 0

	// initial call has empty apfrom
	res := wikibot.GetWikiPageBatch(wikiUrl, "", entryLimit)

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

		res = wikibot.GetWikiPageBatch(wikiUrl, res.Continue.Apcontinue, entryLimit-entries)
	}

	fmt.Printf("📖 Found %d entries \n", entries)
}

// TODO - support for more formats: csv, xdxf, etc.
func (d Dict) Write(path string, format string) string {
	formattedText := f.Format(format, *d.Lex)

	fileName := fmt.Sprintf("%s.%s", d.Name, format)
	normalizedPath := filepath.Join(filepath.FromSlash(path), fileName)

	file, fileErr := os.Create(normalizedPath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	defer file.Close()

	_, writeErr := file.WriteString(formattedText)
	if fileErr != nil {
		log.Fatal(writeErr)
	}

	fmt.Printf("Successfully built dictionary at %s\n", normalizedPath)

	return normalizedPath
}
