package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GetRequestOptions struct {
	// The go http client can't negotiate connections to wikis using TLS 1.2
	ForceTLS12 bool
}

func GetRequest[T any](baseUrl string, queryParams url.Values, options GetRequestOptions) (T, error) {
	queryUrl, err := url.Parse(baseUrl)
	if err != nil {
		return *new(T), err
	}
	queryUrl.RawQuery = queryParams.Encode()

	client := &http.Client{}
	if options.ForceTLS12 {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
			},
		}
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	req.Header.Add("User-Agent", "Runik")

	res, err := client.Do(req)
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
