package dict

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/runik-3/builder/internal/wikitext"
	wikibot "github.com/runik-3/builder/wikiBot"
)

type Entry struct {
	Word       string
	Definition string
}

type Lexicon map[string]Entry

func (l Lexicon) Add(e Entry) {
	l[e.Word] = e
}

func (l Lexicon) Print() {
	fmt.Println("Lexicon (definition -- word)")
	fmt.Println("-------------------------------")
	i := 1
	for _, v := range l {
		fmt.Printf("%d. %s -- %s\n", i, v.Word, v.Definition)
		i++
	}
}

type Dict struct {
	Name string
	Lex  Lexicon
}

func (d Dict) GenerateDefinitionsFromWiki(wikiUrl string, depth int, entryLimit int) {
	entries := 0

	// initial call has empty apfrom
	res := wikibot.GetWikiPageBatch(wikiUrl, "", entryLimit)

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			def := wikitext.ParseDefinition(p.Revisions[0].Slots.Main.Content, depth)
			if def != "" {
				d.Lex.Add(Entry{Word: p.Title, Definition: def})
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
	formattedText := Format(format, d.Lex)

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
