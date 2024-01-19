package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func GetRequest[T any](baseUrl string, queryParams url.Values) (T, error) {
	queryUrl, err := url.Parse(baseUrl)
	if err != nil {
		return *new(T), err
	}
	queryUrl.RawQuery = queryParams.Encode()

	res, err := http.Get(queryUrl.String())
	if err != nil {
		return *new(T), err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var data T

	e := json.Unmarshal(body, &data)
	if e != nil {
		return *new(T), err
	}

	return data, nil
}
