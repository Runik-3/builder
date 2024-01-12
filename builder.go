package main

import (
	d "github.com/runik-3/builder/dict"
	"github.com/runik-3/builder/internal/utils"
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
