package dict

import (
	j "encoding/json"
	"errors"
	"fmt"
)

func Format(format string, l Lexicon) (string, error) {
	switch format {
	case "json":
		j, fmtErr := json(l)
		if fmtErr != nil {
			return "", fmtErr
		}
		return j, nil

	case "df":
		return df(l), nil
	case "csv":
		return csv(l), nil
	case "xdxf":
		return xdxf(l), nil
	default:
		return "", errors.New(fmt.Sprintf("Unsupported file format detected: %s \n", format))
	}
}

func json(l Lexicon) (string, error) {
	json, marshalErr := j.Marshal(l)
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
