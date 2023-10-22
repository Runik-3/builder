package formatter

import (
	j "encoding/json"
	"fmt"
	"log"

	"github.com/runik-3/builder/pkg/lexicon"
)

func Format(format string, l lexicon.Lexicon) string {
	switch format {
	case "json":
		return json(l)
	case "df":
		return df(l)
	case "csv":
		return csv(l)
	case "xdxf":
		return xdxf(l)
	default:
		log.Fatalf("Unsupported file format detected: %s \n", format)
		return ""
	}
}

func json(l lexicon.Lexicon) string {
	json, marshalErr := j.Marshal(l)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	return string(json)
}

func df(l lexicon.Lexicon) string {
	dictFile := ""

	for _, v := range l {
		dictFile += fmt.Sprintf("@ %s\n%s\n", v.Word, v.Definition)
	}

	return dictFile
}

func csv(l lexicon.Lexicon) string {
	// TODO implement
	return ""
}

func xdxf(l lexicon.Lexicon) string {
	// TODO implement
	return ""
}
