package dict

import (
	"fmt"
	"os"
	"path"
	"strings"

	j "encoding/json"

	"github.com/pgaskin/dictutil/dictgen"
	"github.com/pgaskin/dictutil/kobodict"
	_ "github.com/pgaskin/dictutil/kobodict/marisa"
)

// currently only supports Kobo readers
func ConvertForReader(target string, pathToRawDict string, outputDir string) (string, error) {
	// read raw dictionary and unmarshal as Dict
	rawDict, err := os.ReadFile(pathToRawDict)
	if err != nil {
		fmt.Println(err.Error())
	}

	dict := Dict{}
	err = j.Unmarshal(rawDict, &dict)

	df, err := Format("df", dict)
	if err != nil {
		fmt.Println(err.Error())
	}

	dictFile, err := dictgen.ParseDictFile(strings.NewReader(df))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = dictFile.Validate()
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := os.Create(path.Join(outputDir, "dicthtml-"+dict.Name+".zip"))
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	kw := kobodict.NewWriter(file)

	err = dictFile.WriteDictzip(kw, new(dictgen.ImageHandlerRemove), dictgen.ImageFuncFilesystem)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer kw.Close()

	return "", nil
}
