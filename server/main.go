package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var jsonFilePath string
var generateJsonFlag string

func init() {
	flag.StringVar(&jsonFilePath, "file", "", "path")
	flag.StringVar(&generateJsonFlag, "gen-json", "", "generate json from pics folder")
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
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ping")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("endpoint to root")
		for k, v := range imageMap {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	})
	http.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var imageMap map[string]string
		json.Unmarshal(reqBody, &imageMap)
		for k, v := range imageMap {
			fmt.Fprintf(w, "%s: %s\n", k, v)
			askClassifier(v)
		}
	})
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func fileExists(filename string) bool {
	if info, ok := exists(filename); ok {
		return !info.IsDir()
	}
	return false
}

func dirExists(filename string) bool {
	if info, ok := exists(filename); ok {
		return info.IsDir()
	}
	return false
}

func exists(filename string) (fs.FileInfo, bool) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Fatalf("Filename %s not exists", filename)
		return nil, false
	}
	return info, true
}

func askClassifier(filename string) {
	//	if !fileExists(filename) {
	//		return
	//	}

	//TODO PATH?
	cmd := exec.Command("python3", "../classifier/main.py", "--fifo", filename)
	_, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	_, err = cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	//	 panic("foo")
	//
	//		go copyOutput(stdout)
	//		go copyOutput(stderr)
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		//	fmt.Println(scanner.Text())
	}
}

func main() {
	flag.Parse()

	fmt.Printf("jsonFilePath: %v\n", jsonFilePath)

	if generateJsonFlag != "" && dirExists(generateJsonFlag) {
		generateJson(generateJsonFlag, jsonFilePath)
	}

	var imageMap map[string]string
	if jsonFilePath != "" {

		jsonFile, err := os.Open(jsonFilePath)
		defer jsonFile.Close()
		if err != nil {
			panic(err)
		}

		byteJsonValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJsonValue, &imageMap)
	}

	handleRequests(imageMap)
}
