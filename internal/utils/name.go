package utils

import (
	"log"
	"net/url"
	"strings"
)

// pulls name from wiki subdomain
// https://red-rising.fandom.com/api.php ==> red-rising
func NameFromWiki(wikiUrl string) string {
	u, err := url.Parse(wikiUrl)

	if err != nil {
		log.Fatal("please enter a valid wiki url")
	}

	return strings.Split(u.Hostname(), ".")[0]
}
