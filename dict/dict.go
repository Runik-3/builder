package dict

import (
	"fmt"
	"path/filepath"

	"github.com/runik-3/builder/internal/utils"
	"github.com/runik-3/builder/internal/wikitext"
	wikibot "github.com/runik-3/builder/wikiBot"
)

type Entry struct {
	Word       string
	Definition string
}

type Lexicon []Entry

func (l *Lexicon) Add(e Entry) {
	*l = append(*l, e)
}

func (l *Lexicon) Print() {
	fmt.Println("Lexicon (definition -- word)")
	fmt.Println("-------------------------------")
	i := 1
	for _, v := range *l {
		fmt.Printf("%d. %s -- %s\n", i, v.Word, v.Definition)
		i++
	}
}

type Dict struct {
	Name    string
	Lexicon Lexicon
}

func (d *Dict) GenerateDefinitionsFromWiki(wikiUrl string, depth int, entryLimit int) error {
	entries := 0

	// initial call has empty apfrom
	res, batchErr := wikibot.GetWikiPageBatch(wikiUrl, "", entryLimit) //**
	if batchErr != nil {
		return batchErr
	}

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			def := wikitext.ParseDefinition(p.Revisions[0].Slots.Main.Content, depth)
			if def != "" {
				d.Lexicon.Add(Entry{Word: p.Title, Definition: def})
				entries++
			}
		}

		if entries == entryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}

		res, batchErr = wikibot.GetWikiPageBatch(wikiUrl, res.Continue.Apcontinue, entryLimit-entries)
		if batchErr != nil {
			return batchErr
		}
	}

	fmt.Printf("📖 Found %d entries \n", entries)
	return nil
}

// TODO - support for more formats: csv, xdxf, etc.
func (d Dict) Write(path string, format string) (string, error) {
	fmtText, err := Format(format, d)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s.%s", d.Name, format)
	normalizedPath := filepath.Join(filepath.FromSlash(path), fileName)

	err = utils.WriteToFile(fmtText, normalizedPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Successfully built dictionary at %s\n", normalizedPath)

	return normalizedPath, nil
}
