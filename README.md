# Runik Builder

Generate e-reader compatible dictionaries for your favourite fictional worlds.

Builder is the underlying technology at the heart of Runik, containing the tooling required to parse wikis and generate dictionaries.

## Motivation

Runik builder uses the crowd-sourced info of fan wikis to generate dictionaries with custom definitions for fictional names, places, and items. Builder uses non-proprietary formats to store definitions like `json` and `dictfile`, allowing for interop between other tools.

Builder can be used on its own as a CLI tool or as a module import.

## Quick start

```bash
git clone git@github.com:Runik-3/builder.git
cd builder
go mod tidy

go run . info [wiki_url]
```

## CLI use

The quickest and easiest way to get started using Runik Builder is via the command line.

### `info`

Validate a wiki url and retrieve metadata describing the wiki. Takes a wiki url as its argument.

#### Example

```bash
builder info https://red-rising.fandom.com/api.php
```

The following will print to the standard output:

```
Wiki title: Red Rising Wiki
Language: en
Size: 753 entries
Other supported languages:
  - Spanish: https://amanecer-rojo.fandom.com/es/wiki/
  - Hungarian: https://voros-lazadas.fandom.com/hu/wiki/
  - Turkish: http://kizil-yukselis.fandom.com/tr/wiki/
```

### `generate`

Generates an e-reader dictionary based on the pages of a mediawiki-compatible wiki. The generate command requires a wiki url and takes a series of optional flags.

A mediawiki URL used as the target to parse and build the dictionary (eg. https://stardust.fandom.com/api.php).

#### Example

```bash
builder generate https://stardust.fandom.com/api.php
```

#### Options

| Flag | Name          | Description                                                                                                                                                                                                                                                                                                                         |
| ---- | ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-n` | name          | The file name of the generated dictionary file (extension added automatically). If no name is passed in, the file name will default to the subdomain of the target wiki (eg. `red-rising.fandom.com` becomes `red-rising.json`).                                                                                                    |
| `-o` | Out directory | The directory where the generated dictionary will be written. If no directory is specified, a file is not written to disk.                                                                                                                                                                                                          |
| `-f` | Format        | The file format the dictionary is written in. Builder currently supports writing to json and dictfile (`'df'`). When no format is specified, json is the default.                                                                                                                                                                   |
| `-d` | Depth         | The number of sentences that make up each definition. Builder starts parsing the definition at the first setence of the article's main content, ignoring tables and other formatting. As such, be wary that a greater depth has a higher probability of including spoilers. Values between 1 and 5 are suggested, the default is 1. |
| `-l` | Limit         | The maximum number of word entries written to a dictionary. Useful for testing. If no limit is specified, the default is 10,000.                                                                                                                                                                                                    |

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

### Get Wiki Details

The `GetWikiDetails` function is exported from `/wikiBot`. It can act as a tool to validate wiki urls or simply fetch useful metadata about a wiki before it's generated.

```go
func GetWikiDetails(wikiUrl string) (WikiDetails, error)

// returns WikiDetails
type WikiDetails struct {
	SiteName  string
	MainPage  string
	Lang      string
	Logo      string
	Pages     int
	Articles  int
	Languages []Lang
}
```

### Generate Dictionary

Import the `/dict` package to get access to the `BuildDictionary` function, takes a wikiurl and some other options and builds a dictionary based on its entries.

```go
func BuildDictionary(wikiUrl string, options GeneratorOptions) (Dict, error)

// takes GeneratorOptions
type GeneratorOptions struct {
	Name       string
	Output     string
	Depth      int
	Format     string
	EntryLimit int
}

// returns Dict
type Dict struct {
	Name    string
	Lexicon Lexicon
}

type Lexicon []Entry

type Entry struct {
	Word       string
	Definition string
}
```
