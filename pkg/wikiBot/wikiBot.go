package wikibot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
func GetWikiPageBatch(baseUrl string, apfrom string, limit int) *AllPagesResponse {
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

	queryUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	queryUrl.RawQuery = params.Encode()

	res, resErr := http.Get(queryUrl.String())
	if resErr != nil {
		log.Fatalf(resErr.Error())
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var data AllPagesResponse
	e := json.Unmarshal(body, &data)
	if e != nil {
		log.Fatalf(e.Error())
	}

	return &data
}

func PagesToFetch(left int) int {
	// maximum page entries you can fetch and still get full revisions
	const maxPages = 50

	if left < maxPages {
		return left
	}
	return maxPages
}
