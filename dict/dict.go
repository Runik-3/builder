package dict

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

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

func (l *Lexicon) Find(word string) (Entry, bool) {
	for _, entry := range *l {
		if entry.Word == word {
			return entry, true
		}
	}
	return Entry{}, false
}

func (l *Lexicon) Print() {
	fmt.Println("Lexicon (definition -- word)")
	fmt.Println("-------------------------------")
	for i, v := range *l {
		fmt.Printf("%d. %s - %s\n", i+1, v.Word, v.Definition)
	}
}

type Dict struct {
	Name    string
	Lexicon Lexicon
}

type BatchFunction func(src string, startFrom string, limit int) (wikibot.AllPagesResponse, error)

func (d *Dict) GenerateDefinitionsFromWiki(getBatch BatchFunction, wikiUrl string, options GeneratorOptions) error {
	entries := 0
	// initial call starts from the first page
	startFrom := ""

	generateStart := time.Now()
	cont := true
	for cont {
		res, err := getBatch(wikiUrl, startFrom, options.EntryLimit-entries)
		if err != nil {
			return err
		}

		// Break out if we're in an endless loop and not processing anything.
		const TIMEOUT = 15 * time.Second
		elapsed := time.Since(generateStart)
		if elapsed > time.Duration(TIMEOUT) {
			break
		}

		// FIXME: order is not guaranteed in the Lexicon since we're iterating
		// a map of pages.
		for _, p := range res.Query.Pages {
			def, err := wikitext.ParseDefinition(p.Revisions[0].Slots.Main.Content, options.Depth)
			if err != nil {
				continue
			}
			if def != "" {
				d.Lexicon.Add(Entry{Word: p.Title, Definition: def})
				entries++
			}
		}

		if entries == options.EntryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}
		// the next batch call starts on the page where we left off
		startFrom = res.Continue.Apcontinue
	}

	// TODO only print in CLI mode
	fmt.Printf("📖 Found %d entries \n", entries)
	return nil
}

// pulls name from wiki subdomain
// https://red-rising.fandom.com/api.php ==> red-rising
func (d *Dict) NameFromWiki(wikiUrl string) (*Dict, error) {
	u, err := url.Parse(wikiUrl)
	if err != nil {
		return d, errors.New("Must be a valid wiki url")
	}

	dictName := strings.Split(u.Hostname(), ".")[0]
	d.Name = dictName

	return d, nil
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
