package wikibot

import (
	"net/url"
	"strconv"

	"github.com/runik-3/builder/pkg/utils"
)

type AllPagesResponse struct {
	Batchcomplete string   `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         AllPages `json:"query"`
}

type Continue struct {
	Apcontinue string `json:"gapcontinue"`
	Continue   string `json:"continue"`
}

type AllPages struct {
	Pages map[string]Page `json:"pages"`
}

type Page struct {
	Ns        int        `json:"ns"`
	PageId    int        `json:"pageid"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
}

// -- todo make this interface compatible with the deprecated revisions structure
type Revision struct {
	Slots Slot `json:"slots"`
}

type Slot struct {
	Main Main `json:"main"`
}

type Main struct {
	Model   string `json:"contentmodel"`
	Format  string `json:"contentformat"`
	Content string `json:"*"`
}

// fetches batch of entries and unmarshalls the result
func GetWikiPageBatch(baseUrl string, apfrom string, limit int) AllPagesResponse {
	// build query params
	params := url.Values{
		"action":    {"query"},
		"generator": {"allpages"},
		"gaplimit":  {strconv.Itoa(PagesToFetch(limit))},
		"gapfrom":   {apfrom},
		"prop":      {"revisions"},
		"rvprop":    {"content"},
		"rvslots":   {"main"},
		"format":    {"json"},
	}

	return utils.GetRequest[AllPagesResponse](baseUrl, params)
}

func PagesToFetch(left int) int {
	// maximum page entries you can fetch and still get full revisions
	const MaxPages = 50

	if left < MaxPages {
		return left
	}
	return MaxPages
}
