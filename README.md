# Runik Builder

A library that generates e-reader compatible dictionaries for your favourite fictional works using MediaWikis.

Using the content from MediaWiki sites, Builder parses each page as a dictionary entry -- the title of the page becomes the word, and the content becomes the definition. 

## CLI use

The quickest and easiest way to get started using Runik Builder is via the command line.

### Get Wiki Info

Validate a wiki url and retrieve some metadata about the contents of the wiki.

#### Example

```bash
builder info https://red-rising.fandom.com/api.php
```

The following will print to the standard output:

```
Wiki title: Red Rising Wiki
Language: en
Size: 750 entries
Other supported languages: Spanish, Hungarian, Turkish
```

### Generate Dictionary

Generates an e-reader dictionary based on the pages of a mediawiki-compatible wiki. The generate command requires a wiki url and takes a series of optional flags.

A mediawiki URL used as the target to parse and build the dictionary (eg. https://stardust.fandom.com/api.php).

#### Example

```bash
builder generate https://stardust.fandom.com/api.php
```

#### Flags

`-n` Name

The file name of the generated dictionary file (extension added automatically). If no name is passed in, the file name will default to the subdomain of the target wiki (eg. `red-rising.fandom.com` becomes `red-rising.json`).

`-o` Output directory

The directory where the generated dictionary will be written. If no directory is specified, a file is not written to disk.

`-f` Format

The file format the dictionary is written in. Builder currently supports writing to json and dictfile (`'df'`). When no format is specified, json is the default.

`-d` Depth

The number of sentences that make up each definition. Be wary that a greater depth has a higher probability of including spoilers. The default is 1.

`-l` Limit

The maximum number of word entries written to a dictionary. Useful for testing. If no limit is specified, the default is 10,000.

#### Example

```
builder -o ./ -f df generate https://stardust.fandom.com/api.php
```

Running this command generates a dictionary from the entire Stardust Fandom Wiki and write its contents to the current directory as `stardust.df`.

```
...
@ Wall Guard
The Wall Guard was a 97-year-old man who guarded the gap in the wall which was the border between Stormhold and England.
@ Ingrid
Ingrid is the star that fell 400 years before Yvaine did, more precisely at some point (presumably) in the 15th century.
...
```

A snippet of the resulting dictfile.

## Module use

Builder can be imported as a module to use in your own projects. Run the following command in the root of your project to add builder.

```bash
go get github.com/runik-3/builder
```

### Get Wiki Info

The `GetWikiDetails` function is exported from `/wikibot`. It can act as a tool to validate wiki urls or simply fetch useful metadata about a wiki before it's generated.

```go
type WikiDetails struct {
	SiteName  string
	MainPage  string
	Lang      string
	Logo      string
	Pages     int
	Articles  int
	Languages []Lang
}

func GetWikiDetails(wikiUrl string) WikiDetails
```

### Generate Dictionary

Import the `/dict` package to get access to the `BuildDictionary` function, takes a wikiurl and some other options and builds a dictionary based on its entries.

```go
type Dict struct {
	Name string
	Lex  Lexicon
}

type Lexicon map[string]Entry

type Entry struct {
	Word       string
	Definition string
}

func BuildDictionary(
  wikiUrl string,
  name string,
  output string,
  entryLimit int,
  depth int,
  format string
) Dict
```
