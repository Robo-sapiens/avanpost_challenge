package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var jsonFilePath string

func init() {
  flag.StringVar(&jsonFilePath, "file", "", "path")
}

func main() {

  flag.Parse()

  fmt.Printf("jsonFilePath: %v\n", jsonFilePath)

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
