package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetRequest[T any](baseUrl string, queryParams url.Values) T {
	queryUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	queryUrl.RawQuery = queryParams.Encode()

	res, resErr := http.Get(queryUrl.String())
	if resErr != nil {
		log.Fatalf(resErr.Error())
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var data T

	e := json.Unmarshal(body, &data)
	if e != nil {
		log.Fatalf(e.Error())
	}

	return data
}
