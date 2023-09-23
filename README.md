# Runik Builder

A library that generates e-reader compatible dictionaries from MediaWikis. 

## Quick Start

CLI flags: 
- `-w`: `[*required]` Must be a valid MediaWiki api url (eg. https://stardust.fandom.com/api.php). 
- `-l`: `[default 10,000]` Sets a limit for 

The following command will generate a dictionary from the Stardust Fandom Wiki containing the first 5 words and their definitions.
```
`go run . -w https://stardust.fandom.com/api.php -l 5`
```
