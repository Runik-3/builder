package utils

import (
	"errors"
	"net/url"
	"strings"
)

// ensures the url points to the wiki's api endpoint
func NormalizeUrl(u string, suffix string) (string, error) {
	parsedUrl, err := url.Parse(u)

	if strings.Contains(parsedUrl.Path, suffix) {
		return u, nil
	}

	invalidErr := errors.New("Invalid url, please try again with a valid url (eg. https://malazan.fandom.com/api.php)")
	if err != nil {
		return "", invalidErr
	}

	if parsedUrl.Host == "" || parsedUrl.Scheme == "" {
		return "", invalidErr
	}

	endpointUrl := url.URL{}

	endpointUrl.Scheme = parsedUrl.Scheme
	endpointUrl.Host = parsedUrl.Host
	endpointUrl.Path = buildWikiPath(parsedUrl.Path, suffix)

	return endpointUrl.String(), nil

}

// builds a wiki url path, retaining a language code if the wiki is not
// english language
func buildWikiPath(path string, suffix string) string {
	pathParts := strings.Split(path, "/")
	builtPath := ""

	if len(pathParts) <= 1 {
		return suffix
	}

	startOfUrlPath := pathParts[1]

	// check if the start of the url path is a valid iso language code
	for _, langCode := range languageCodes {
		if len(pathParts) > 1 && langCode == startOfUrlPath {
			builtPath = "/" + langCode
		}
	}

	return builtPath + suffix
}

var languageCodes = []string{
	"ab",
	"aa",
	"af",
	"sq",
	"am",
	"ar",
	"hy",
	"as",
	"ay",
	"az",
	"ba",
	"eu",
	"bn",
	"dz",
	"bh",
	"bi",
	"br",
	"bg",
	"my",
	"be",
	"km",
	"ca",
	"zh",
	"zh",
	"co",
	"hr",
	"cs",
	"da",
	"nl",
	"en",
	"eo",
	"et",
	"fo",
	"fa",
	"fj",
	"fi",
	"fr",
	"fy",
	"gl",
	"gd",
	"gv",
	"ka",
	"de",
	"el",
	"kl",
	"gn",
	"gu",
	"ha",
	"he",
	"hi",
	"hu",
	"is",
	"id",
	"ia",
	"ie",
	"iu",
	"ik",
	"ga",
	"it",
	"ja",
	"ja",
	"kn",
	"ks",
	"kk",
	"rw",
	"ky",
	"rn",
	"ko",
	"ku",
	"lo",
	"la",
	"lv",
	"li",
	"ln",
	"lt",
	"mk",
	"mg",
	"ms",
	"ml",
	"mt",
	"mi",
	"mr",
	"mo",
	"mn",
	"na",
	"ne",
	"no",
	"oc",
	"or",
	"om",
	"ps",
	"pl",
	"pt",
	"pa",
	"qu",
	"rm",
	"ro",
	"ru",
	"sm",
	"sg",
	"sa",
	"sr",
	"sh",
	"st",
	"tn",
	"sn",
	"sd",
	"si",
	"ss",
	"sk",
	"sl",
	"so",
	"es",
	"su",
	"sw",
	"sv",
	"tl",
	"tg",
	"ta",
	"tt",
	"te",
	"th",
	"bo",
	"ti",
	"to",
	"ts",
	"tr",
	"tk",
	"tw",
	"ug",
	"uk",
	"ur",
	"uz",
	"vi",
	"vo",
	"cy",
	"wo",
	"xh",
	"yi",
	"yo",
	"zu",
}
