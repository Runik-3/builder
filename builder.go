package main

import (
	"fmt"

	d "github.com/runik-3/builder/dict"
	"github.com/runik-3/builder/internal/utils"
	wikibot "github.com/runik-3/builder/wikiBot"
)

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, depth int, format string) d.Dict {
	dict := d.Dict{Lex: d.Lexicon{}}

	dictName := name
	if dictName == "" {
		dictName = utils.NameFromWiki(wikiUrl)
	}
	dict.Name = dictName

	dict.GenerateDefinitionsFromWiki(utils.FormatUrl(wikiUrl), depth, entryLimit)

	if output != "" {
		dict.Write(output, format)
	}

	return dict
}

func printWikiDetails(wikiUrl string) {
	details := wikibot.GetWikiDetails(wikiUrl)

	fmt.Printf("Wiki title: %s\n", details.SiteName)
	fmt.Printf("Language: %s\n", details.Lang)
	fmt.Printf("Size: %d entries\n", details.Articles)

	langs := []wikibot.Lang{}
	for _, lang := range details.Languages {
		langs = append(langs, lang)
	}

	if len(langs) > 0 {
		fmt.Print("Other supported languages: ")
		// supported langs
		for i, lang := range langs {
			fmt.Printf("%s", lang.LangName)
			if i < len(langs)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
}
