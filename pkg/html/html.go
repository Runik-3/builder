package html

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func ParseFirstPara(doc string) (string, error) {
	r := strings.NewReader(doc)
	t := html.NewTokenizer(r)

	var prevToken html.Token
	str := ""
	firstP := false

	for {
		tt := t.Next()
		tkn := t.Token()

		switch tt {

		// end of file or error
		case html.ErrorToken:
			if t.Err() == io.EOF {
				return "", errors.New("Could not find intro paragraph in docs")
			}

		case html.StartTagToken:
			if tkn.Data == "b" && prevToken.Data == "p" {
				firstP = true
			}
			prevToken = tkn

		case html.EndTagToken:
			if firstP && tkn.Data == "p" {
				return str, nil
			}

		case html.TextToken:
			if firstP {
				str += tkn.Data
			}
			prevToken = tkn
		}
	}
}
