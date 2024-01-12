package main

import (
	"flag"
)

func main() {
	wikiUrl := flag.String("w", "", "The wiki api url (eg. https://examplewiki.org/api.php).")
	entryLimit := flag.Int("l", 10000, "The maximum number of entries in the dictionary.")
	depth := flag.Int("d", 1, "How many sentences make up each definition.")
	name := flag.String("n", "", "The file name of the generated dictionary (extension added automatically).")
	output := flag.String("o", "", "The output directory where generated dictionary will be written. If not passed in, no file will be generated (preferred behaviour when calling builder as module).")
	format := flag.String("f", "json", "Format of the output file.")

	flag.Parse()

	BuildDictionary(*wikiUrl, *name, *output, *entryLimit, *depth, *format)
}
