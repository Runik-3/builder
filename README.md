# Runik Builder

This library generates e-reader compatible dictionaries from MediaWikis. 

## Quick Start

The command line currently accepts two flags: 
- `-u`: Is the wiki api url (eg. https://stardust.fandom.com/api.php).
- `-p`: Is the page limit (defaults to 5 for development purposes). <-- currently this controls the size of mediawiki page queries, but this will get swapped for a total dictionary size limiter.


Example: 
```
`go run . -u https://stardust.fandom.com/api.php -p 5`
```
