package dict

import (
	"bytes"
	j "encoding/json"
	"fmt"
	"text/template"
)

// Formats dictionaries into raw text
func Format(format string, d Dict) (string, error) {
	switch format {
	case "json":
		j, fmtErr := json(d)
		if fmtErr != nil {
			return "", fmtErr
		}
		return j, nil

	case "df":
		dictFile, err := df(d)
		return dictFile, err
	default:
		return "", fmt.Errorf("Unsupported file format detected: %s \n", format)
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

func df(d Dict) (string, error) {
	const DF_TEMPLATE = `{{range .Lexicon}}@ {{.Word}}
{{range .Synonyms}}& {{.}}
{{end}}{{.Definition}}
{{end}}`
	tmpl, err := template.New("DictFile template").Parse(DF_TEMPLATE)
	if err != nil {
		return "", nil
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, d)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
