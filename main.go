package main

import (
	"flag"

	d "github.com/runik-3/builder/pkg/dict"
	l "github.com/runik-3/builder/pkg/lexicon"
	"github.com/runik-3/builder/pkg/utils"
)

func main() {
	wikiUrl := flag.String("w", "", "The wiki api url (eg. https://examplewiki.org/api.php).")
	entryLimit := flag.Int("l", 10000, "The maximum number of entries in the dictionary.")
	name := flag.String("n", "", "The file name of the generated dictionary (extension added automatically).")
	output := flag.String("o", "", "The output directory where generated dictionary will be written. If not passed in, no file will be generated (preferred behaviour when calling builder as module).")
	format := flag.String("f", "json", "Format of the output file.")

	flag.Parse()

	BuildDictionary(*wikiUrl, *name, *output, *entryLimit, *format)
}

func BuildDictionary(wikiUrl string, name string, output string, entryLimit int, format string) d.Dict {
	dictName := name
	if dictName == "" {
		dictName = utils.NameFromWiki(wikiUrl)
	}

	dict := d.Dict{Name: dictName, Lex: l.New()}
	dict.GenerateDefinitionsFromWiki(wikiUrl, entryLimit)
	if output != "" {
		dict.Write(output, format)
	}

	return dict
}
