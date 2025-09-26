package wikiBot

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

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
	ApiUrl    string
	SiteName  string
	MainPage  string
	Lang      string
	Logo      string
	Pages     int
	Articles  int
	Languages []Lang
	// options to include in requests to the wiki server
	RequestOpts utils.GetRequestOptions
}

// fetches details about the requested wiki including,
// name, size, and supported languages.
// returns err if wiki url is invalid.
func GetWikiDetails(wikiUrl string) (WikiDetails, error) { // do some parsing and return a better shape
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("meta", "siteinfo")
	params.Add("siprop", "statistics|general")
	params.Add("origin", "*")

	requestOpts := utils.GetRequestOptions{
		// Is is this timeout too short?
		Timeout: 5 * time.Second,
	}

	// Mediawiki api endpoints can have any of these suffixes
	apiSuffixes := []string{"/api.php", "/w/api.php", "/wiki/api.php"}
	wikiApiEndpoint := ""
	var wikiDetailsRes WikiDetailsResponse

	// Try potential api endpoint suffixes
	for _, suffix := range apiSuffixes {
		apiUrl, err := utils.NormalizeUrl(wikiUrl, suffix)
		if err != nil {
			continue
		}

		wikiDetailsRes, err = utils.GetRequest[WikiDetailsResponse](apiUrl, params, requestOpts)
		// wikis using tls 1.2 can throw 403s because the go http module can't
		// establish a connection. Repeat the same request, but force the request to
		// use tls 1.2 to be compatible.
		if err != nil && strings.Contains(err.Error(), "403") {
			requestOpts = utils.GetRequestOptions{ForceTLS12: true}
			wikiDetailsRes, err = utils.GetRequest[WikiDetailsResponse](apiUrl, params, requestOpts)
		}
		if err == nil {
			wikiApiEndpoint = apiUrl
			break
		}
	}

	// TODO-MM: is this sufficient to know if we found the endpoint?
	if wikiApiEndpoint == "" {
		return WikiDetails{}, errors.New("Could not find a valid endpoint for this wiki. Try copy/pasting the exact wiki api endpoint url.")
	}

	wikiLangsRes, err := wikiLanguages(wikiApiEndpoint, wikiDetailsRes.Query.General.MainPage, requestOpts)
	if err != nil {
		return WikiDetails{}, err
	}

	langs := []Lang{}
	for _, p := range wikiLangsRes.Query.Pages {
		langs = append(langs, p.LangLinks...)
	}

	return WikiDetails{
		ApiUrl:      wikiApiEndpoint,
		SiteName:    wikiDetailsRes.Query.General.SiteName,
		MainPage:    wikiDetailsRes.Query.General.MainPage,
		Lang:        wikiDetailsRes.Query.General.Lang,
		Logo:        wikiDetailsRes.Query.General.Logo,
		Pages:       wikiDetailsRes.Query.Statistics.Pages,
		Articles:    wikiDetailsRes.Query.Statistics.Articles,
		Languages:   langs,
		RequestOpts: requestOpts,
	}, nil
}

func wikiLanguages(wikiUrl string, mainPage string, options utils.GetRequestOptions) (WikiDetailsResponse, error) {
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "langlinks")
	params.Add("llprop", "url|langname|autonym")
	params.Add("titles", mainPage)

	langRes, err := utils.GetRequest[WikiDetailsResponse](wikiUrl, params, options)
	if err != nil {
		return WikiDetailsResponse{}, err
	}
	return langRes, nil
}

func PrintWikiDetails(wikiUrl string) error {
	details, err := GetWikiDetails(wikiUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Wiki title: %s\n", details.SiteName)
	fmt.Printf("Language: %s\n", details.Lang)
	fmt.Printf("Size: %d entries\n", details.Articles)

	langs := []Lang{}
	for _, lang := range details.Languages {
		langs = append(langs, lang)
	}

	if len(langs) > 0 {
		fmt.Println("Other supported languages: ")
		// supported langs
		for _, lang := range langs {
			fmt.Printf("  - %s: %s\n", lang.LangName, lang.Url)
		}
	}
	return nil
}
