package dict

import (
	j "encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/runik-3/builder/wikiBot"
)

var mockBatchCalled int = 0

func TestGenerateDefinitionsFromWiki(t *testing.T) {
	t.Run("fetches a single batch when limit is 1", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(MockWikiBatchFunction, "", GeneratorOptions{Depth: 2, EntryLimit: 1})

		isEqual(t, dict.Lexicon[0].Word, "batch1_page1", "")
		contains(t, dict.Lexicon[0].Definition, "This is the first page of the first batch.")
		isEqual(t, mockBatchCalled, 1, "")
		isEqual(t, len(dict.Lexicon), 1, "")

		mockBatchCalled = 0
	})

	t.Run("fetches definitions in batches from all wiki pages", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(MockWikiBatchFunction, "", GeneratorOptions{Depth: 2})

		// Note: order is not guaranteed in the Lexicon since we're iterating a
		// map of pages
		def, exists := dict.Lexicon.Find("batch1_page1")
		isEqual(t, exists, true, "")
		contains(t, def.Definition, "This is the first page of the first batch.")

		def, exists = dict.Lexicon.Find("batch2_page1")
		isEqual(t, exists, true, "")
		contains(t, def.Definition, "This is the first page of the second batch.")

		def, exists = dict.Lexicon.Find("batch2_page2")
		isEqual(t, exists, true, "")
		contains(t, def.Definition, "This is the second page of the second batch.")

		def, exists = dict.Lexicon.Find("batch3_page1")
		isEqual(t, exists, true, "")
		contains(t, def.Definition, "This is the first page of the third batch.")

		isEqual(t, mockBatchCalled, 3, "")
		isEqual(t, len(dict.Lexicon), 4, "")

		mockBatchCalled = 0
	})
}

// in this mock src determines the kind of response we get back
func MockWikiBatchFunction(src string, startFrom string, limit int) (wikiBot.AllPagesResponse, error) {
	mockBatchCalled += 1

	allBatches := getFixtureData("wikiResponseBatch.json")
	batch := wikiBot.AllPagesResponse{}

	// simulate fetching batches based on start page from the API
	batch = allBatches[startFrom]

	return batch, nil
}

func getFixtureData(fixture string) map[string]wikiBot.AllPagesResponse {
	pathToResponseJson, _ := filepath.Abs(filepath.Join("fixtures", fixture))
	responseJson, _ := os.ReadFile(pathToResponseJson)

	var response map[string]wikiBot.AllPagesResponse
	j.Unmarshal(responseJson, &response)

	return response
}

func isEqual(t *testing.T, val1 any, val2 any, message string) {
	if val1 != val2 {
		if message != "" {
			t.Fatal(message)
		} else {
			t.Fatalf("%v not equal to %v", val1, val2)
		}
	}
}

func contains(t *testing.T, str string, subStr string) {
	if !strings.Contains(str, subStr) {
		t.Fatalf("%v does not contain %v", str, subStr)
	}
}
