package dict

import (
	j "encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/runik-3/builder/wikiBot"
)

var mockBatchCalled int = 0

func TestGenerateDefinitionsFromWiki(t *testing.T) {
	t.Run("Gets a single batch of definitions", func(t *testing.T) {
		dict := Dict{}
		dict.GenerateDefinitionsFromWiki(MockWikiBatchFunction, "", GeneratorOptions{Depth: 2})

		if dict.Lexicon[0].Word != "10th Legion" {
			t.Fatal("Did not properly parse dictionary word")
		}
		fmt.Println(dict.Lexicon[0].Definition)
		if !strings.Contains(dict.Lexicon[0].Definition, "The 10th Legion was one of the three legions in the Malaz") {
			t.Fatal("Did not properly parse dictionary definition")
		}
		// TODO mockBatchCalled should only get called once when fetching a single batch
		if mockBatchCalled > 2 {
			t.Fatalf("Batch function called %d times", mockBatchCalled)
		}
	})
}

// helpers
func MockWikiBatchFunction(src string, startFrom string, limit int) (wikiBot.AllPagesResponse, error) {
	mockBatchCalled += 1
	response := getFixtureData("allPagesResponse.json")

	response.Continue.Apcontinue = ""

	return response, nil
}

func getFixtureData(fixture string) wikiBot.AllPagesResponse {
	pathToResponseJson, _ := filepath.Abs(filepath.Join("fixtures", fixture))
	responseJson, _ := os.ReadFile(pathToResponseJson)

	var response wikiBot.AllPagesResponse
	j.Unmarshal(responseJson, &response)

	return response
}
