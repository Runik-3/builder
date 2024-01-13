package wikibot

import (
	"fmt"
	"net/url"

	"github.com/runik-3/builder/internal/utils"
)

type WikiDetailsResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         Query  `json:"query"`
}

type Query struct {
	Pages      map[string]Page `json:"pages"`
	Statistics Statistics      `json:"statistics"`
	General    General         `json:"general"`
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

type WikiDetails struct {
	SiteName  string
	MainPage  string
	Lang      string
	Logo      string
	Pages     int
	Articles  int
	Languages []Lang
}

// fetches details about the requested wiki including,
// name, size, and supported languages.
// returns err if wiki url is invalid.
func GetWikiDetails(wikiUrl string) WikiDetails { // do some parsing and return a better shape
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("meta", "siteinfo")
	params.Add("siprop", "statistics|general")

	details := utils.GetRequest[WikiDetailsResponse](wikiUrl, params)
	langsDetails := wikiLanguages(wikiUrl, details.Query.General.MainPage)

	langs := []Lang{}
	for _, p := range langsDetails.Query.Pages {
		langs = append(langs, p.LangLinks...)
	}

	return WikiDetails{
		SiteName:  details.Query.General.SiteName,
		MainPage:  details.Query.General.MainPage,
		Lang:      details.Query.General.Lang,
		Logo:      details.Query.General.Logo,
		Pages:     details.Query.Statistics.Pages,
		Articles:  details.Query.Statistics.Articles,
		Languages: langs,
	}
}

func wikiLanguages(wikiUrl string, mainPage string) WikiDetailsResponse {
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "langlinks")
	params.Add("llprop", "url|langname|autonym")
	params.Add("titles", mainPage)

	return utils.GetRequest[WikiDetailsResponse](wikiUrl, params)
}

func PrintWikiDetails(wikiUrl string) {
	details := GetWikiDetails(wikiUrl)

	fmt.Printf("Wiki title: %s\n", details.SiteName)
	fmt.Printf("Language: %s\n", details.Lang)
	fmt.Printf("Size: %d entries\n", details.Articles)

	langs := []Lang{}
	for _, lang := range details.Languages {
		langs = append(langs, lang)
	}

	if len(langs) > 0 {
		fmt.Print("Other supported languages: ")
		// supported langs
		for i, lang := range langs {
			fmt.Printf("%s", lang.LangName)
			if i < len(langs)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
}
