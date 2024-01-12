package wikibot

import (
	"net/url"

	"github.com/runik-3/builder/pkg/utils"
)

type WikiDetailsResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         Query  `json:"query"`
}

type Query struct {
	Pages      Page       `json:"pages"`
	Statistics Statistics `json:"statistics"`
	General    General    `json:"general"`
}

type Statistics struct {
	Pages    int `json:"pages"`
	Articles int `json:"articles"`
}

type General struct {
	MainPage string `json:"mainpage"`
	SiteName string `json:"sitename"`
	Logo     string `json:"logo"`
	Lang     string `json:"lang"`
}

// fetches details about the requested wiki including,
// name, size, and supported languages.
// returns err if wiki url is invalid.
func WikiDetails(wikiUrl string) WikiDetailsResponse {
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "langlinks")
	params.Add("meta", "siteinfo")
	params.Add("llprop", "url|langname|autonym")
	params.Add("siprop", "statistics|general")

	return utils.GetRequest[WikiDetailsResponse](wikiUrl, params)
}
