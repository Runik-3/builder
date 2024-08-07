package wikiBot

import (
	"net/url"
	"strconv"

	"github.com/runik-3/builder/internal/utils"
)

type AllPagesResponse struct {
	Batchcomplete string   `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         Query    `json:"query"`
}

type Continue struct {
	Apcontinue string `json:"gapcontinue"`
	Continue   string `json:"continue"`
}

type Pages struct {
	Pages map[string]Page `json:"pages"`
}

type Page struct {
	PageId    int        `json:"pageid"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
	LangLinks []Lang     `json:"langlinks"`
}

type Lang struct {
	Lang     string `json:"lang"`
	LangName string `json:"langname"`
	Autonym  string `json:"autonym"`
	Url      string `json:"url"`
}

// TODO make this interface compatible with the deprecated revisions structure
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
func GetWikiPageBatch(baseUrl string, startFrom string, limit int) (AllPagesResponse, error) {
	// define query params
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("generator", "allpages")
	params.Add("gaplimit", strconv.Itoa(pagesToFetch(limit)))
	params.Add("gapfrom", startFrom)
	params.Add("prop", "revisions")
	params.Add("rvprop", "content")
	params.Add("rvslots", "main")

	res, err := utils.GetRequest[AllPagesResponse](baseUrl, params)
	if err != nil {
		return AllPagesResponse{}, err
	}
	return res, nil
}

func pagesToFetch(left int) int {
	// maximum page entries you can fetch and still get full revisions
	const MaxPages = 50

	if left < MaxPages {
		return left
	}
	return MaxPages
}
