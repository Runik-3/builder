# Runik Builder

A library that generates e-reader compatible dictionaries from MediaWikis.

## Quick Start

### CLI use

#### `-w` Wiki target

The mediawiki URL used as a target to build the dictionary (eg. https://stardust.fandom.com/api.php).

#### `-n` Name

The file [n]ame of the generated dictionary file (extension added automatically). If no name is passed in, the file name will default to the subdomain of the target wiki (eg. red-rising.fandom.com becomes red-rising.json).

#### `-o` Output directory

The directory where the generated dictionary will be written. If no directory is specified, a file is not written to disk.

#### `-f` Format

The file format the dictionary is written in. Builder currently supports writing to json and dictfile (`'df'`). When no format is specified, json is the default.

#### `-l` Limit

The maximum number of word entries written to a dictionary. Useful for testing. If no limit is specified, the default is 10,000.

#### Example

The command below will generate a dictionary from the entire Stardust Fandom Wiki and write its contents to the current directory as `stardust.df`.

```
go run . -w https://stardust.fandom.com/api.php -o ./ -f df
```
