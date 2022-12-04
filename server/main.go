package main

import (
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
	"sync"
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

	cmd := "python3"
	args := []string{"../classifier/main.py", "--cli"}

	tCmd := exec.Command(cmd, args...)

	stdin, err := tCmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "430__F_Right_little_finger_CR.BMP")
	}()

  tCmd.Wait()

	out, err := tCmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
  fmt.Println("sfdg")

	fmt.Printf("string(out): %v\n", string(out))
	tasks := make(chan *exec.Cmd, 64)

	var wg sync.WaitGroup
	// env var this size TODO
	sizeOf := 4
	for k, v := range imageMap {
		for i := 0; i < sizeOf; i++ {
			wg.Add(1)
			go func(num int, w *sync.WaitGroup) {
				defer w.Done()

				var (
					out   []byte
					err   error
					stdin io.WriteCloser
				)

				for cmd := range tasks {
					fmt.Println("into tasks")
					stdin, err = cmd.StdinPipe()
					if err != nil {
						log.Fatalln(err)
					}
					defer stdin.Close()
					fmt.Printf("v: %v\n", v)
					io.WriteString(stdin, v)

					cmd.Wait()
					out, err = cmd.Output()
					fmt.Printf("cmd.Dir: %v\n", cmd.Dir)
					fmt.Printf("cmd.Err: %v\n", cmd.Err)
					fmt.Printf("cmd.Path: %v\n", cmd.Path)
					fmt.Printf("cmd.Args: %v\n", cmd.Args)
					if err != nil {
						log.Fatal("can't get stdout: ", err)
					}
					fmt.Printf("goroutine %d commands output %s\n", num, string(out))
				}
			}(i, &wg)
		}
		fmt.Printf("k: %v\n", k)
	}

	tasks <- exec.Command(cmd, args...)

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
		tmpMap := make(map[string]string, 1)
		for k, v := range imageMap {
			tmpMap[k] = v
			break
		}
		fmt.Printf("tmpMap: %v\n", tmpMap)
		askClassifier(tmpMap)
	}

	if isWeb {
		handleRequests(imageMap)
	}
}
