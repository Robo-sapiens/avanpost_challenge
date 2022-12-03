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
	"sync"
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
      if !fileExists(v) {
        log.Printf("File %s with filepath %s do not exists\n", k, v)
      }
		}
		askClassifier(imageMap)
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

func askClassifier(imageMap map[string]string) {

	cmd := "python 3 ../classifier/main.py"
	args := []string{"--fifo"}

	tasks := make(chan *exec.Cmd, 64)

	var wg sync.WaitGroup
	// env var this size TODO
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(num int, w *sync.WaitGroup) {
			defer w.Done()

			var (
				out []byte
				err error
			)

			for cmd := range tasks {
				out, err = cmd.Output()
				if err != nil {
					log.Fatal("can't get stdout:", err)
				}
				fmt.Printf("goroutine %d commands output %s\n", num, string(out))
			}
		}(i, &wg)
	}

	//TODO Process file ?
	for k, v := range imageMap {
		args = append(args, v)
		tasks <- exec.Command(cmd, args...)
		fmt.Printf("k: %v\n", k)
	}

	close(tasks)

	wg.Wait()

	fmt.Println("end")

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
