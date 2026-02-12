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

func (d *Dict) GenerateDefinitionsFromWiki(getBatch BatchFunction, wiki wikibot.WikiDetails, options GeneratorOptions) error {
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

		for _, p := range res.Query.Pages {
			// We didn't retrieve the full batch of redirects. We have the same page
			// content, with some additional redirects to add to exising entries.
			if redirectsContinue != "" {
				if len(p.Redirects) == 0 {
					continue
				}

				existingEntry, found := d.Lexicon.Find(p.Title)
				if !found {
					continue
				}
				for _, r := range p.Redirects {
					existingEntry.Synonyms = append(existingEntry.Synonyms, r.Title)
				}
			} else {
				// Parse page content
				word := wikitext.ParseWord(p.Title)
				def, err := wikitext.ParseDefinition(p.GetPageContent(), options.Depth)
				if err != nil || def == "" {
					continue
				}

				redirects := []string{}
				if len(p.Redirects) > 0 {
					for _, r := range p.Redirects {
						redirects = append(redirects, r.Title)
					}
				}

				d.Lexicon.Add(Entry{Word: word, Definition: def, Synonyms: redirects})
				entries++
			}
		}

		if options.ProgressHook != nil {
			total := wiki.Articles
			if options.EntryLimit < wiki.Articles {
				total = options.EntryLimit
			}
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
			// reset
			redirectsContinue = ""
			// the next batch call starts on the page where we left off
			startFrom = res.Continue.Gapcontinue
		}
	}

	// TODO only print in CLI mode
	fmt.Printf("📖 Found %d entries \n", entries)
	return nil
}

// TODO - support for more formats: csv, xdxf, etc.
func (d Dict) Write(path string, format string) (string, error) {
	d.Lexicon.Sort()

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
