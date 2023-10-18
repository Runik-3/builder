package formatter

import (
	j "encoding/json"
	"log"

	"github.com/runik-3/builder/pkg/lexicon"
)

func Format(format string, l lexicon.Lexicon) string {
	switch format {
	case "json":
		return json(l)
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

func csv(l lexicon.Lexicon) string {
	// TODO implement
	return ""
}

func xdxf(l lexicon.Lexicon) string {
	// TODO implement
	return ""
}
