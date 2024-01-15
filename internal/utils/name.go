package utils

import (
	"errors"
	"fmt"
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
	fmt.Println(wikiUrl, u.Hostname())

	return strings.Split(u.Hostname(), ".")[0], nil
}
