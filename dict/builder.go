package dict

import (
	"github.com/runik-3/builder/internal/utils"
	wikibot "github.com/runik-3/builder/wikiBot"
)

type GeneratorOptions struct {
	Name       string
	Output     string
	Depth      int
	Format     string
	EntryLimit int
}

func BuildDictionary(wikiUrl string, options GeneratorOptions) (Dict, error) {
	dict := Dict{Lexicon: Lexicon{}}

	if options.Name != "" {
		dict.Name = options.Name
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

	err = dict.GenerateDefinitionsFromWiki(wikibot.GetWikiPageBatch, u, options)
	if err != nil {
		return Dict{}, err
	}

	if options.Output != "" {
		dict.Write(options.Output, options.Format)
	}

	return dict, nil
}
