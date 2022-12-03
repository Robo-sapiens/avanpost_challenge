package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func handleRequests(imageMap map[string]string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("endpoint to root")
		for k, v := range imageMap {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	})
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var imageMap map[string]string
		json.Unmarshal(reqBody, &imageMap)
		// TODO UNIMPLEMENT
		// after get result?     json.NewEncoder(w).Encode(...)
	})
	log.Fatal(http.ListenAndServe(":10000", nil))
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

	handleRequests(imageMap)
}
