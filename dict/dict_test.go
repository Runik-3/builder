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

		if dict.Lexicon[0].Word != "batch1_page1" {
			t.Fatal("Did not properly parse dictionary word")
		}
		if !strings.Contains(dict.Lexicon[0].Definition, "This is the first page of the first batch.") {
			t.Fatal("Did not properly parse dictionary definition")
		}
		if mockBatchCalled != 1 {
			t.Fatalf("Batch function called %d times", mockBatchCalled)
		}
		if len(dict.Lexicon) != 1 {
			t.Fatalf("Dict contains %d entries", mockBatchCalled)
		}
		mockBatchCalled = 0
	})

	t.Run("fetches definitions for all pages of a wiki", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(MockWikiBatchFunction, "", GeneratorOptions{Depth: 2})

		if dict.Lexicon[0].Word != "batch1_page1" {
			t.Fatal("Did not properly parse dictionary word")
		}
		if !strings.Contains(dict.Lexicon[0].Definition, "This is the first page of the first batch.") {
			t.Fatal("Did not properly parse dictionary definition")
		}
		if mockBatchCalled != 3 {
			t.Fatalf("Batch function called %d times", mockBatchCalled)
		}
		if len(dict.Lexicon) != 4 {
			t.Fatalf("Dict contains %d entries", mockBatchCalled)
		}
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
