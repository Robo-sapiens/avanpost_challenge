package main

import (
	"bytes"
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
	//
	//	tCmd := exec.Command(cmd, args...)
	//
	//	stdin, err := tCmd.StdinPipe()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	go func() {
	//		defer stdin.Close()
	//		io.WriteString(stdin, "430__F_Right_little_finger_CR.BMP")
	//	}()
	//
	//  tCmd.Wait()
	//
	//	out, err := tCmd.CombinedOutput()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//  fmt.Printf("tCmd.Dir: %v\n", tCmd.Dir)
	//  fmt.Printf("tCmd.Args: %v\n", tCmd.Args)
	//  fmt.Printf("tCmd.Path: %v\n", tCmd.Path)
	//  fmt.Printf("tCmd.Env: %v\n", tCmd.Env)
	//  fmt.Println("sfdg")
	//
	//	fmt.Printf("string(out): %v\n", string(out))
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

					//	cmd.Wait()
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

//	fmt.Printf("jsonFilePath: %v\n", jsonFilePath)

	if generateJsonFlag != "" && dirExists(generateJsonFlag) {
		generateJson(generateJsonFlag, jsonFilePath)
	}
//	foo("51__M_Left_ring_finger_CR.BMP")

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
//		tmpMap := make(map[string]string, 1)
//		for k, v := range imageMap {
//			tmpMap[k] = v
//      fmt.Printf("%v\n", v)
//			break
//		}
//		askClassifier(tmpMap)
	}

	if isWeb {
		handleRequests(imageMap)
	}
}

// CapturingPassThroughWriter is a writer that remembers
// data written to it and passes it to w
type CapturingPassThroughWriter struct {
	buf bytes.Buffer
	w   io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func foo(file string) {

	py := "python3"
	args := []string{"../classifier/main.py", "--cli"}
	cmd := exec.Command(py, args...)
  

	var errStdout, errStderr error
//	stdoutIn, _ := cmd.StdoutPipe()
//	stderrIn, _ := cmd.StderrPipe()
	stdout := NewCapturingPassThroughWriter(os.Stdout)
	stderr := NewCapturingPassThroughWriter(os.Stderr)
	stdin := NewCapturingPassThroughWriter(os.Stdin)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		//_, errStdout = io.Copy(stdout, stdoutIn)
    io.WriteString(stdin, file)
		wg.Done()
	}()

//	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
  out, err :=  cmd.Output()
  if err != nil {
    log.Fatalf("output error is %s", err)
  }
  fmt.Printf("string(out): %v\n", string(out))
	if errStdout != nil || errStderr != nil {
		log.Fatalf("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
