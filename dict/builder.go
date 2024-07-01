package dict

import (
	"github.com/runik-3/builder/internal/utils"
)

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, depth int, format string) (Dict, error) {
	dict := Dict{Lexicon: Lexicon{}}

	if name != "" {
		dict.Name = name
	} else {
		_, err := dict.NameFromWiki(wikiUrl)
		if err != nil {
			return Dict{}, err
		}
	}

	u, err := utils.NormalizeUrl(wikiUrl)
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
