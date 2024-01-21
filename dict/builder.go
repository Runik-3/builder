package dict

import (
	"github.com/runik-3/builder/internal/utils"
)

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, depth int, format string) (Dict, error) {
	dict := Dict{Lexicon: Lexicon{}}

	dictName := name
	if dictName == "" {
		nameFromWiki, err := utils.NameFromWiki(wikiUrl)
		if err != nil {
			return Dict{}, err
		}
		dictName = nameFromWiki
	}
	dict.Name = dictName

	u, err := utils.FormatUrl(wikiUrl)
	if err != nil {
		return Dict{}, err
	}

	err = dict.GenerateDefinitionsFromWiki(u, depth, entryLimit)
	if err != nil {
		return Dict{}, err
	}

	if output != "" {
		dict.Write(output, format)
	}

	return dict, nil
}
