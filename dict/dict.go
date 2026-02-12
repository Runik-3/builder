package dict

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/runik-3/builder/internal/utils"
	"github.com/runik-3/builder/internal/wikitext"
	wikibot "github.com/runik-3/builder/wikiBot"
)

type Entry struct {
	Word       string
	Definition string
	Synonyms   []string
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

func (l *Lexicon) Sort() {
	sort.Slice(*l, func(i, j int) bool {
		return strings.ToLower((*l)[i].Word) < strings.ToLower((*l)[j].Word)
	})
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
	ApiUrl  string
	Lang    string
	Lexicon Lexicon
}

type BatchFunction func(src string, startFrom string, limit int, redirectsContinue string, options utils.GetRequestOptions) (wikibot.AllPagesResponse, error)

func (dict *Dict) GenerateDefinitionsFromWiki(getBatch BatchFunction, wiki wikibot.WikiDetails, options GeneratorOptions) error {
	entries := 0
	startFrom := ""
	redirectsContinue := ""

	cont := true
	for cont {
		res, err := getBatch(wiki.ApiUrl, startFrom, options.EntryLimit-entries, redirectsContinue, wiki.RequestOpts)
		if err != nil {
			return err
		}

		if len(res.Query.Pages) == 0 {
			return errors.New("Could not get page content.")
		}

		for _, page := range res.Query.Pages {
			if redirectsContinue != "" {
				handleParseAdditionalRedirects(page, dict)
			} else {
				entry, ok := parseContentAsEntry(page, options)
				if !ok {
					continue
				}

				dict.Lexicon.Add(entry)
				entries++
			}
		}

		if options.ProgressHook != nil {
			total := min(options.EntryLimit, wiki.Articles)
			options.ProgressHook(entries, total)
		}

		if entries == options.EntryLimit {
			break
		}

		if res.Continue.Gapcontinue == "" && res.Continue.Rdcontinue == "" {
			cont = false
		}

		if res.Continue.Rdcontinue != "" {
			// we didn't get all the redirects for this batch, continue fetching
			redirectsContinue = res.Continue.Rdcontinue
		} else {
			// the next batch call starts on the page where we left off
			startFrom = res.Continue.Gapcontinue
			// reset before next request
			redirectsContinue = ""
		}
	}

	// TODO only print in CLI mode
	fmt.Printf("📖 Found %d entries \n", entries)
	return nil
}

// TODO - support for more formats: csv, xdxf, etc.
func (dict Dict) Write(path string, format string) (string, error) {
	dict.Lexicon.Sort()

	fmtText, err := Format(format, dict)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s.%s", dict.Name, format)
	normalizedPath := filepath.Join(filepath.FromSlash(path), fileName)

	err = utils.WriteToFile(fmtText, normalizedPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Successfully built dictionary at %s\n", normalizedPath)

	return normalizedPath, nil
}

func parseContentAsEntry(page wikibot.Page, options GeneratorOptions) (Entry, bool) {
	word := wikitext.ParseWord(page.Title)
	def, err := wikitext.ParseDefinition(page.GetPageContent(), options.Depth)
	if err != nil || def == "" {
		return Entry{}, false
	}

	redirects := []string{}
	if len(page.Redirects) > 0 {
		for _, r := range page.Redirects {
			redirects = append(redirects, r.Title)
		}
	}

	return Entry{Word: word, Definition: def, Synonyms: redirects}, true
}

// We didn't retrieve the full batch of redirects. We have the same page
// content, with some additional redirects to add to exising entries.
func handleParseAdditionalRedirects(p wikibot.Page, d *Dict) {
	if len(p.Redirects) == 0 {
		return
	}

	existingEntry, found := d.Lexicon.Find(p.Title)
	if !found {
		return
	}

	for _, r := range p.Redirects {
		existingEntry.Synonyms = append(existingEntry.Synonyms, r.Title)
	}
}
