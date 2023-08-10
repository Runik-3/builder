package wikibot

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"cgt.name/pkg/go-mwclient"
)

type Response struct {
	Batchcomplete bool     `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         AllPages `json:"query"`
}

type AllPages struct {
	Pages []Page `json:"allpages"`
}

type Page struct {
	Ns     int    `json:"ns"`
	PageId int    `json:"pageid"`
	Title  string `json:"title"`
}

type Continue struct {
	Apcontinue string `json:"apcontinue"`
	Continue   string `json:"continue"`
}

// all below funcs belong in here, but perhaps this method gets brought out into
// it's own pkg?
func GenerateDictionary() {
	wikiUrl := flag.String("u", "", "wikiUrl")
	// low default limit of 5 for testing, should be 500 in prod
	pageLimit := flag.Int("p", 5, "pageLimit")
	flag.Parse()

	w := CreateClient(*wikiUrl)
	GetWikiPages(w, pageLimit, "")
}

func CreateClient(url string) *mwclient.Client {
	// format for fandom == <property_name>.fandom.com/api.php
	w, err := mwclient.New(url, "myWikibot")
	if err != nil {
		panic(err)
	}

	return w
}

func GetWikiPages(w *mwclient.Client, pageLimit *int, apfrom string) Response {
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

	var data Response
	json.Unmarshal([]byte(resp), &data)

	fmt.Printf("📖 Found %d entries \n", len(data.Query.Pages))

	return data
}
