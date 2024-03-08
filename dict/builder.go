package dict

import (
	"github.com/runik-3/builder/internal/utils"
)

type GeneratorOptions struct {
	Name       string
	Output     string
	EntryLimit int
	Depth      int
	Format     string
}

func BuildDictionary(wikiUrl string, options GeneratorOptions) (Dict, error) {
	dict := Dict{Lexicon: Lexicon{}}

	dictName := options.Name
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

	err = dict.GenerateDefinitionsFromWiki(u, options.Depth, options.EntryLimit)
	if err != nil {
		return Dict{}, err
	}

	if options.Output != "" {
		dict.Write(options.Output, options.Format)
	}

	return dict, nil
}
