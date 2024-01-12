package utils

import (
	"log"
	"net/url"
	"strings"
)

// ensures the url points to the wiki's api endpoint
func FormatUrl(u string) (validUrl string) {
	parsedUrl, e := url.Parse(u)
	if e != nil {
		log.Fatalf("Invalid url, please try again with a valid url (eg. https://malazan.fandom.com/api.php)")
	}

	if parsedUrl.Host == "" || parsedUrl.Scheme == "" {
		log.Fatalf("Invalid url, please try again with a valid url (eg. https://malazan.fandom.com/api.php)")
	}

	endpointUrl := url.URL{}

	if !strings.Contains(parsedUrl.Path, "api.php") {
		endpointUrl.Scheme = parsedUrl.Scheme
		endpointUrl.Host = parsedUrl.Host
		endpointUrl.Path = "/api.php"
		return endpointUrl.String()
	}

	return u
}
