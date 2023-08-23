package wikibot

import (
	"encoding/json"
	"fmt"
	"github.com/runik-3/builder/pkg/dict"
	"strconv"

	"cgt.name/pkg/go-mwclient"
)

type AllPagesRes struct {
	Batchcomplete bool     `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         AllPages `json:"query"`
}

type Continue struct {
	Apcontinue string `json:"apcontinue"`
	Continue   string `json:"continue"`
}

type AllPages struct {
	Pages []Page `json:"allpages"`
}

type Page struct {
	Ns     int    `json:"ns"`
	PageId int    `json:"pageid"`
	Title  string `json:"title"`
}

// perhaps this function gets brought out into it's own pkg?
// also this should take args so it can be used as a module
// main can check if flags weren't passed and panic
func GenerateWordList(d *dict.Dict, wikiUrl *string, pageLimit *int) {
	w := CreateClient(*wikiUrl)

	// initial call has empty apfrom
	res := GetWikiPages(w, pageLimit, "")

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			d.Add(dict.Entry{Word: p.Title, Definition: ""})
		}

		if res.Continue.Apcontinue != "" {
			res = GetWikiPages(w, pageLimit, res.Continue.Apcontinue)
		} else {
			cont = false
		}
	}
}

func CreateClient(url string) *mwclient.Client {
	w, err := mwclient.New(url, "myWikibot")
	if err != nil {
		panic(err)
	}

	return w
}

func GetWikiPages(w *mwclient.Client, pageLimit *int, apfrom string) *AllPagesRes {
	parameters := map[string]string{
		"action":  "query",
		"list":    "allpages",
		"aplimit": strconv.Itoa(*pageLimit),
		"apfrom":  apfrom,
	}

	resp, err := w.GetRaw(parameters)
	if err != nil {
		panic(err)
	}

	var data AllPagesRes
	json.Unmarshal([]byte(resp), &data)

	fmt.Printf("📖 Found %d entries between %s and %s \n", len(data.Query.Pages), data.Query.Pages[0].Title, data.Query.Pages[len(data.Query.Pages)-1].Title)

	return &data
}
