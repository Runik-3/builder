package builder

import (
	d "github.com/runik-3/builder/pkg/dict"
	"github.com/runik-3/builder/pkg/utils"
)

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, depth int, format string) {
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
}
