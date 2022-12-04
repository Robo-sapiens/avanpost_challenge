package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
var isWeb bool

func init() {
	flag.StringVar(&jsonFilePath, "file", "", "path")
	flag.StringVar(&generateJsonFlag, "gen-json", "", "generate json from pics folder")
	flag.BoolVar(&isWeb, "web", false, "web server run")
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

type Response struct {
  FileNum string `json:"fileNum"`
  FileName string `json:"fileName"`
  PyResponse string `json:"classifierResponse"`
}

type ResponseJSON []Response

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
    var responses ResponseJSON
		for k, v := range imageMap {
			fmt.Fprintf(w, "%s: %s\n", k, v)
			if !fileExists(v) {
				log.Printf("File %s with filepath %s do not exists\n", k, v)
			}
      responses = append(responses, Response{
        FileNum: k,
        FileName: v,
        PyResponse: processImage(v),
      })
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
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
	ge := os.Getenv("PICS")
	info, err := os.Stat(ge + filename)
	if os.IsNotExist(err) {
		log.Printf("Filename %s not exists", filename)
		return nil, false
	}
	return info, true
}

func main() {
	flag.Parse()

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
		for _, v := range imageMap {
			fmt.Printf("%v\n", v)
		}
	}

	if isWeb {
		handleRequests(imageMap)
	}
}

func processImage(file string) string {
	py := "python3"
	args := []string{"../classifier/main.py", "--image_file", file}
	subprocess := exec.Command(py, args...)

	out, err := subprocess.Output()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("string(out): %v", string(out))
	subprocess.Wait()
	return string(out)
}
