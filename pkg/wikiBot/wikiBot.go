package wikibot

import (
	"encoding/json"
	"fmt"
	"github.com/runik-3/builder/pkg/dict"

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

func GetWikiPages(w *mwclient.Client, apfrom string) *AllPagesRes {
	params := map[string]string{
		"action":  "query",
		"list":    "allpages",
		"aplimit": "500",
		"apfrom":  apfrom,
	}

	resp, err := w.GetRaw(params)
	if err != nil {
		panic(err)
	}

	var data AllPagesRes
	json.Unmarshal([]byte(resp), &data)
	println(len(data.Query.Pages))

	return &data
}

//func ParsePage(w *mwclient.Client, id string) {
//	params := map[string]string{
//		"action":  "query",
//		"list":    "allpages",
//	}
//}

func GenerateWordList(d *dict.Dict, wikiUrl *string, entryLimit *int) {
	w := CreateClient(*wikiUrl)
	// Should notify the user when a limit is reached and there is more to the dict
	entries := 0

	// initial call has empty apfrom
	res := GetWikiPages(w, "")

	// continue?
	cont := true

	for cont {
		for _, p := range res.Query.Pages {
			if entries == *entryLimit {
				fmt.Printf("📖 Found %d entries \n", entries)
				return
			}

			d.Add(dict.Entry{Word: p.Title, Definition: ""})
			entries++
		}

		if res.Continue.Apcontinue == "" {
			cont = false
		}
		// TODO - future optimization check entry limit to only fetch required data using aplimit
		res = GetWikiPages(w, res.Continue.Apcontinue)
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
