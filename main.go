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
	name := flag.String("n", "", "The file name of the generated dictionary.")
	// output
	// definition depth (how many lines)
	// print
	flag.Parse()

	BuildDictionary(*wikiUrl, *name, *entryLimit)
}

func BuildDictionary(wikiUrl string, name string, entryLimit int) {
	dictName := name
	if dictName == "" {
		dictName = utils.NameFromWiki(wikiUrl)
	}

	dict := d.Dict{Name: dictName, Lex: l.New()}
	dict.GenerateDefinitionsFromWiki(wikiUrl, entryLimit)
	dict.Lex.Print()
	dict.Write("")
}
