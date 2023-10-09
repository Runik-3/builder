# Runik Builder

A library that generates e-reader compatible dictionaries from MediaWikis. 

## Quick Start

CLI flags: 
- `-w`: `[*required]` Must be a valid Media[w]iki api url (eg. https://stardust.fandom.com/api.php). 
- `-n`: `[defaults to wiki subdomain (eg. red-rising.fandom.com becomes red-rising.json)]` Custom dictionary [n]ame.
- `-o`: `[by default no file is written]` The [o]utput directory where the generated dictionary will be written.
- `-l`: `[default 10,000]` Sets a [l]imit for 

The following command will generate a dictionary from the Stardust Fandom Wiki containing the first 5 words and their definitions.
```
go run . -w https://stardust.fandom.com/api.php -l 5
```
