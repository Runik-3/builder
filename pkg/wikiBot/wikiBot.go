package wikibot

import (
	"dictGen/pkg/dict"
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

// perhaps this function gets brought out into it's own pkg?
// also this should take args so it can be used as a module
// main can check if flags weren't passed and panic
func GenerateWordList() {
	wikiUrl := flag.String("u", "", "wikiUrl")
	// low default limit of 5 for testing, should be 500
	pageLimit := flag.Int("p", 5, "pageLimit")
	flag.Parse()

	w := CreateClient(*wikiUrl)

	// initial call has empty apfrom
	res := GetWikiPages(w, pageLimit, "")
	fmt.Println(res)

	d := dict.New()
	d.Add(dict.Entry{Word: "Test", Definition: "Sup"})
	d.Print()
}

func CreateClient(url string) *mwclient.Client {
	w, err := mwclient.New(url, "myWikibot")
	if err != nil {
		panic(err)
	}

	return w
}

func GetWikiPages(w *mwclient.Client, pageLimit *int, apfrom string) *Response {
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
	fmt.Println(data)

	return &data
}
