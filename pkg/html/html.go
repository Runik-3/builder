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

	inP := false
	firstP := false
	content := ""

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
			if tkn.Data == "p" {
				inP = true
			}
			if tkn.Data == "b" && inP {
				firstP = true
			} else {
				content = ""
			}

		case html.EndTagToken:
			if firstP && tkn.Data == "p" {
				return content, nil
			}

		case html.TextToken:
			if inP || firstP {
				content += tkn.Data
			}
		}
	}
}
