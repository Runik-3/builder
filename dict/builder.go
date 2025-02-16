package dict

import (
	wikibot "github.com/runik-3/builder/wikiBot"
)

type GeneratorOptions struct {
	Name       string
	Output     string
	Depth      int
	Format     string
	EntryLimit int
	// A function that gets called after each batch is processed, used to hook
	// into the progress of the generator.
	ProgressHook func(processed int, total int)
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

	// TODO: We have more info we can use when we want to extend the dictionary
	// json schema to include more context.
	wiki, err := wikibot.GetWikiDetails(wikiUrl)
	if err != nil {
		return Dict{}, err
	}

	err = dict.GenerateDefinitionsFromWiki(wikibot.GetWikiPageBatch, wiki, options)
	if err != nil {
		return Dict{}, err
	}

	if options.Output != "" {
		dict.Write(options.Output, options.Format)
	}

	return dict, nil
}
