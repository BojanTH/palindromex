package service

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const staticURLPrefix = "/static/"
const staticFilesDir = "web/static/dist/"

var staticFilesMap, _ = getStaticFilesMap()

func getStaticFilesMap() (map[string]string, error) {
	filesMap := make(map[string]string)

	err := filepath.Walk(staticFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filenameWithHash := filepath.Base(path)

			hashStart := strings.Index(filenameWithHash, ".")
			hashEnd := strings.LastIndex(filenameWithHash, ".")
			fileName := filenameWithHash[:hashStart] + filenameWithHash[hashEnd:]

			filesMap[staticURLPrefix + fileName] = staticURLPrefix + filenameWithHash
		}

		return nil
	})

	if err != nil {
		return filesMap, err
	}

	return filesMap, nil
}

func getFunctionsMap() template.FuncMap {
    return template.FuncMap{
		// Static function is used for css and js files that have file versioning
		"static": func(fileName string) string {
			if nameWithHash, ok := staticFilesMap[fileName]; ok {
				return nameWithHash
			}
			log.Fatal(fmt.Sprintf("File \"%s\" doesn't exist", fileName))

			return fileName
		},
	}
}
