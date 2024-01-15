package dict

import (
	"github.com/runik-3/builder/internal/utils"
)

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, depth int, format string) (Dict, error) {
	dict := Dict{Lex: Lexicon{}}

	dictName := name
	if dictName == "" {
		nameFromWiki, nameErr := utils.NameFromWiki(wikiUrl)
		if nameErr != nil {
			return Dict{}, nameErr
		}
		dictName = nameFromWiki
	}
	dict.Name = dictName

	genErr := dict.GenerateDefinitionsFromWiki(utils.FormatUrl(wikiUrl), depth, entryLimit)
	if genErr != nil {
		return Dict{}, genErr
	}

	if output != "" {
		dict.Write(output, format)
	}

	return dict, nil
}
