package utils

import (
	"errors"
	"net/url"
	"strings"
)

// pulls name from wiki subdomain
// https://red-rising.fandom.com/api.php ==> red-rising
func NameFromWiki(wikiUrl string) (string, error) {
	u, err := url.Parse(wikiUrl)
	if err != nil {
		return "", errors.New("please enter a valid wiki url")
	}

	return strings.Split(u.Hostname(), ".")[0], nil
}
