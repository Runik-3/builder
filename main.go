package main

import (
	"flag"
	"log"

	"github.com/runik-3/builder/dict"
	wikibot "github.com/runik-3/builder/wikiBot"
)

func main() {
	entryLimit := flag.Int("l", 10000, "The maximum number of entries in the dictionary.")
	depth := flag.Int("d", 1, "How many sentences make up each definition.")
	name := flag.String("n", "", "The file name of the generated dictionary (extension added automatically).")
	output := flag.String("o", "", "The output directory where generated dictionary will be written. If not passed in, no file will be generated (preferred behaviour when calling builder as module).")
	format := flag.String("f", "json", "Format of the output file.")

	flag.Parse()
	args := flag.Args()
	command := args[0]
	wikiUrl := args[1]

	if len(args) < 2 {
		log.Fatal("You must provide at least two arguments.")
	}

	switch command {
	case "generate":

		dict.BuildDictionary(wikiUrl, *name, *output, *entryLimit, *depth, *format)
	case "info":
		wikibot.PrintWikiDetails(wikiUrl)
	default:
		log.Fatalf("%s is not a valid command. See help for more options.", args[0])
	}
}
