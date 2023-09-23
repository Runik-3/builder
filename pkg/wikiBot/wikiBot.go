package wikibot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/runik-3/builder/pkg/dict"

	"cgt.name/pkg/go-mwclient"
)

type AllPagesResponse struct {
	Batchcomplete bool     `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         AllPages `json:"query"`
}

type Continue struct {
	Apcontinue string `json:"apcontinue"`
	Continue   string `json:"continue"`
}

type AllPages struct {
	Pages []Page `json:"pages"`
}

type Page struct {
	Ns        int        `json:"ns"`
	PageId    int        `json:"pageid"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
}

type Revision struct {
	Slots Slot `json:"slots"`
}

type Slot struct {
	Main Main `json:"main"`
}

type Main struct {
	Model   string `json:"contentmodel"`
	Format  string `json:"contentformat"`
	Content string `json:"content"`
}

func GetWikiPages(w *mwclient.Client, apfrom string, limit int) *AllPagesResponse {
	params := map[string]string{
		"action":    "query",
		"generator": "allpages",
		"gaplimit":  strconv.Itoa(limit),
		"gapfrom":   apfrom,
		"prop":      "revisions",
		"rvprop":    "content",
		"rvslots":   "main",
	}

	resp, err := w.GetRaw(params)
	if err != nil {
		panic(err)
	}

	var data AllPagesResponse
	json.Unmarshal([]byte(resp), &data)
	println(string(resp))

	return &data
}

func GenerateWordList(d *dict.Dict, wikiUrl *string, entryLimit *int) {
	w := CreateClient(*wikiUrl)
	// Should notify the user when a limit is reached and there is more to the dict
	entries := 0

	// initial call has empty apfrom
	res := GetWikiPages(w, "", *entryLimit)

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			d.Add(dict.Entry{Word: p.Title, Definition: p.Revisions[0].Slots.Main.Content})
			entries++
		}

		if entries == *entryLimit {
			break
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}

		res = GetWikiPages(w, res.Continue.Apcontinue, *entryLimit-entries)
	}
	fmt.Printf("📖 Found %d entries \n", entries)
}

func CreateClient(url string) *mwclient.Client {
	w, err := mwclient.New(url, "myWikibot")
	if err != nil {
		panic(err)
	}

	return w
}
