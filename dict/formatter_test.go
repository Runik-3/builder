package dict

import (
	"os"
	"path/filepath"
	"testing"

	test "github.com/runik-3/builder/internal/testUtils"
)

var dict Dict = Dict{
	Name:    "Test Dict",
	ApiUrl:  "example.com",
	Lang:    "en",
	Lexicon: Lexicon{
		Entry{Word: "Simple", Definition: "This is a simple entry"},
		Entry{Word: "Simple Two", Definition: "Another simple entry"},
		Entry{
			Word: "Synonyms",
			Definition: "This is an entry with synonyms",
			Synonyms: []string{"First", "Second", "Third"},
		},
	},
}

func TestFormatDictFile(t *testing.T) {
	result, err := Format("df", dict)
	if err != nil {
		t.Fatalf("Failed to format dictFile \n%s\n", err)
	}

	test.IsEqual(t, result, getFixtureData("format_dictfile.df"), "")
}

func getFixtureData(fixture string) string {
	pathToResponseJson, _ := filepath.Abs(filepath.Join("fixtures", fixture))
	text, _ := os.ReadFile(pathToResponseJson)

	return string(text)
}
