package dict

import (
	j "encoding/json"
	"os"
	"path/filepath"
	"testing"

	test "github.com/runik-3/builder/internal/testUtils"
	"github.com/runik-3/builder/internal/utils"
	"github.com/runik-3/builder/wikiBot"
)

var mockBatchCalled int = 0

func TestGenerateDefinitionsFromWiki(t *testing.T) {
	t.Run("fetches a single batch when limit is 1", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(mockWikiBatchFunction, wikiBot.WikiDetails{}, GeneratorOptions{Depth: 2, EntryLimit: 1})

		test.IsEqual(t, dict.Lexicon[0].Word, "batch1_page1", "")
		test.Contains(t, dict.Lexicon[0].Definition, "This is the first page of the first batch.")
		test.IsEqual(t, mockBatchCalled, 1, "")
		test.IsEqual(t, len(dict.Lexicon), 1, "")

		mockBatchCalled = 0
	})

	t.Run("fetches definitions in batches from all wiki pages", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(mockWikiBatchFunction, wikiBot.WikiDetails{}, GeneratorOptions{Depth: 2})

		// Note: order is not guaranteed in the Lexicon since we're iterating a
		// map of pages
		def, exists := dict.Lexicon.Find("batch1_page1")
		test.IsEqual(t, exists, true, "")
		test.Contains(t, def.Definition, "This is the first page of the first batch.")

		def, exists = dict.Lexicon.Find("batch2_page1")
		test.IsEqual(t, exists, true, "")
		test.Contains(t, def.Definition, "This is the first page of the second batch.")

		def, exists = dict.Lexicon.Find("batch2_page2")
		test.IsEqual(t, exists, true, "")
		test.Contains(t, def.Definition, "This is the second page of the second batch.")

		def, exists = dict.Lexicon.Find("batch3_page1")
		test.IsEqual(t, exists, true, "")
		test.Contains(t, def.Definition, "This is the first page of the third batch.")

		test.IsEqual(t, mockBatchCalled, 3, "")
		test.IsEqual(t, len(dict.Lexicon), 4, "")

		mockBatchCalled = 0
	})

	t.Run("parses redirects as synonyms", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(mockWikiBatchFunction, wikiBot.WikiDetails{}, GeneratorOptions{Depth: 2, EntryLimit: 1})

		def, _ := dict.Lexicon.Find("batch1_page1")
		test.IsEqual(t, def.Word, "batch1_page1", "")
		expected_syns := []string{"one", "first"}
		test.IsEqual(t, def.Synonyms[0], expected_syns[0], "")
		test.IsEqual(t, def.Synonyms[1], expected_syns[1], "")
	})

	t.Run("refetches redirects using rdcontinue", func(t *testing.T) {
		dict := Dict{}
		if mockBatchCalled < 10 {
			dict.GenerateDefinitionsFromWiki(MockWikiBatchFuncWithRedirect, wikiBot.WikiDetails{}, GeneratorOptions{Depth: 2, EntryLimit: 1})
		}

		def, _ := dict.Lexicon.Find("batch1_page1")
		test.IsEqual(t, def.Word, "batch1_page1", "")
		expected_syns := []string{"one", "first", "b1_p2"}

		for i := range def.Synonyms {
			test.IsEqual(t, def.Synonyms[i], expected_syns[i], "")
		}
	})
}

func mockWikiBatchFunction(src string, startFrom string, limit int, redirectsContinue string, options utils.GetRequestOptions) (wikiBot.AllPagesResponse, error) {
	mockBatchCalled += 1

	allBatches := getFixturePagesResponse("wikiResponseBatch.json")
	batch := wikiBot.AllPagesResponse{}

	// simulate fetching batches based on start page from the API
	batch = allBatches[startFrom]

	return batch, nil
}

func MockWikiBatchFuncWithRedirect(src string, startFrom string, limit int, redirectsContinue string, options utils.GetRequestOptions) (wikiBot.AllPagesResponse, error) {
	mockBatchCalled += 1

	allBatches := getFixturePagesResponse("wikiResponseBatchWithRedirects.json")
	batch := wikiBot.AllPagesResponse{}

	// simulate fetching batches/redirects based on start page from the API
	cont := startFrom
	if redirectsContinue != "" {
		cont = redirectsContinue
	}
	batch = allBatches[cont]

	return batch, nil
}

func getFixturePagesResponse(fixture string) map[string]wikiBot.AllPagesResponse {
	pathToResponseJson, _ := filepath.Abs(filepath.Join("fixtures", fixture))
	responseJson, _ := os.ReadFile(pathToResponseJson)

	var response map[string]wikiBot.AllPagesResponse
	j.Unmarshal(responseJson, &response)

	return response
}
