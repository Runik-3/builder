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
	// Generator continue -- used when we're fetching another set of pages
	Gapcontinue string `json:"gapcontinue"`
	// Redirect continue -- used when we're fetching more redirects
	Rdcontinue string `json:"rdcontinue"`
	Continue   string `json:"continue"`
}

type Pages struct {
	Pages map[string]Page `json:"pages"`
}

type Page struct {
	PageId    int        `json:"pageid"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
	Redirects []Redirect `json:"redirects"`
	LangLinks []Lang     `json:"langlinks"`
}

func (p *Page) GetPageContent() string {
	if p.Revisions[0].Slots.Main.Content != "" {
		return p.Revisions[0].Slots.Main.Content
	}

	// Try path from the old revisions wiki response
	return p.Revisions[0].Content
}

type Lang struct {
	Lang     string `json:"lang"`
	LangName string `json:"langname"`
	Autonym  string `json:"autonym"`
	Url      string `json:"url"`
}

type Revision struct {
	Slots Slot `json:"slots"`

	// Support the old mediawiki revisions api format.
	Model   string `json:"contentmodel"`
	Format  string `json:"contentformat"`
	Content string `json:"*"`
}

type Slot struct {
	Main Main `json:"main"`
}

type Redirect struct {
	PageId int    `json:"pageid"`
	Title  string `json:"title"`
}

type Main struct {
	Model   string `json:"contentmodel"`
	Format  string `json:"contentformat"`
	Content string `json:"*"`
}

// fetches batch of entries and unmarshalls the result
func GetWikiPageBatch(baseUrl string, startFrom string, limit int, redirectsContinue string, options utils.GetRequestOptions) (AllPagesResponse, error) {
	// define query params
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("generator", "allpages")
	params.Add("gaplimit", strconv.Itoa(pagesToFetch(limit)))
	params.Add("gapfrom", startFrom)
	params.Add("prop", "redirects|revisions")
	params.Add("rvprop", "content")
	params.Add("rvslots", "main")
	params.Add("rdprop", "pageid|title")
	params.Add("rdlimit", "max")
	if redirectsContinue != "" {
		params.Add("rdcontinue", redirectsContinue)
	}

	res, err := utils.GetRequest[AllPagesResponse](baseUrl, params, options)
	// successful response
	if err == nil {
		return res, nil
	}

	return AllPagesResponse{}, err
}

func pagesToFetch(left int) int {
	// maximum page entries you can fetch and still get full revisions
	const MaxPages = 50

	if left < MaxPages {
		return left
	}
	return MaxPages
}
