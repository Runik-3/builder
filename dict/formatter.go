package dict

import (
	j "encoding/json"
	"fmt"
	"log"
)

func Format(format string, l Lexicon) string {
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

func json(l Lexicon) string {
	json, marshalErr := j.Marshal(l)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	return string(json)
}

func df(l Lexicon) string {
	dictFile := ""

	for _, v := range l {
		dictFile += fmt.Sprintf("@ %s\n%s\n", v.Word, v.Definition)
	}

	return dictFile
}

func csv(l Lexicon) string {
	// TODO implement
	return ""
}

func xdxf(l Lexicon) string {
	// TODO implement
	return ""
}
