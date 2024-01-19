package utils

import (
	"errors"
	"net/url"
	"strings"
)

// ensures the url points to the wiki's api endpoint
func FormatUrl(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", errors.New("Invalid url, please try again with a valid url (eg. https://malazan.fandom.com/api.php)")
	}

	if parsedUrl.Host == "" || parsedUrl.Scheme == "" {
		return "", errors.New("Invalid url, please try again with a valid url (eg. https://malazan.fandom.com/api.php)")
	}

	endpointUrl := url.URL{}

	if !strings.Contains(parsedUrl.Path, "api.php") {
		endpointUrl.Scheme = parsedUrl.Scheme
		endpointUrl.Host = parsedUrl.Host
		endpointUrl.Path = "/api.php"
		return endpointUrl.String(), nil
	}

	return u, nil
}
