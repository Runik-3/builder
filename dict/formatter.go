package dict

import (
	j "encoding/json"
	"errors"
	"fmt"
)

// Formats dictionaries into raw text
func Format(format string, d Dict) (string, error) {
	lex := d.Lexicon

	switch format {
	case "json":
		j, fmtErr := json(d)
		if fmtErr != nil {
			return "", fmtErr
		}
		return j, nil

	case "df":
		return df(lex), nil
	case "csv":
		return csv(lex), nil
	case "xdxf":
		return xdxf(lex), nil
	default:
		return "", errors.New(fmt.Sprintf("Unsupported file format detected: %s \n", format))
	}
}

// Default format for storing Runik dictionaries. Includes metadata as well as
// definitions.
func json(d Dict) (string, error) {
	json, marshalErr := j.Marshal(d)
	if marshalErr != nil {
		return "", marshalErr
	}

	return string(json), nil
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
