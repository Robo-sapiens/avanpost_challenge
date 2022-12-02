package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var jsonFilePath string
var generateJsonFlag bool

func init() {
	flag.StringVar(&jsonFilePath, "file", "", "path")
	flag.BoolVar(&generateJsonFlag, "gen-json", false, "generate json")
}

func generateJson(relPath string, input string) {
	absPath, _ := filepath.Abs(relPath)

	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		log.Fatal(err)
	}

  jsonData := map[string]string{}
	for i, file := range files {
		fieldName := fmt.Sprintf("image_%d", i)
		jsonData[fieldName] = file.Name()
	}

	jsonFile, _ := json.MarshalIndent(jsonData, "  ", "")
	ioutil.WriteFile(input, jsonFile, 0644)
}

func main() {

	flag.Parse()

	fmt.Printf("jsonFilePath: %v\n", jsonFilePath)

	if generateJsonFlag {
		generateJson("SOCOFing/Altered/Altered-Easy", jsonFilePath)
	}

	jsonFile, err := os.Open(jsonFilePath)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}

	byteJsonValue, _ := ioutil.ReadAll(jsonFile)

	var imageMap map[string]string

	json.Unmarshal(byteJsonValue, &imageMap)

	keys := make([]string, 0, len(imageMap))
	for k := range imageMap {
		keys = append(keys, k)
	}

	for _, v := range keys {
		fmt.Printf("v: %v\n", v)
	}
}
