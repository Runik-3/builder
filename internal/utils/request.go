package utils

import (
	"encoding/json"
	"errors"
	"fmt"
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

	statusOk := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOk {
		return *new(T), errors.New(fmt.Sprintf("Request failed with %s", res.Status))
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
