package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

var knownHashes map[string]string
var soFar = 0

func recursivelyScanDirectory(logBase string, pathBase string, path string, maptouse map[string]string, logfd *os.File) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(files); i++ {

		logMsg := logBase + "> " + strconv.Itoa(i + 1) + "/" + strconv.Itoa(len(files)) + " "
		fmt.Print("\u001B[2K\r", logMsg, " ( ", soFar, " duplicates found )")

		if files[i].IsDir() {
			recursivelyScanDirectory(logMsg, pathBase + files[i].Name() + "/", path + "/" + files[i].Name(), maptouse, logfd)
		} else {
			file, err := os.Open(path + "/" + files[i].Name())
			if err != nil {
				panic(err)
			}

			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				panic(err)
			}

			result := hex.EncodeToString(hash.Sum(nil))

			file.Close()

			if _, ok := maptouse[result]; ok == true {
				soFar++
				logfd.WriteString(pathBase + files[i].Name())
				logfd.WriteString("\n")
				logfd.WriteString(maptouse[result])
				logfd.WriteString("\n======\n")
			} else {
				maptouse[result] = pathBase + files[i].Name()
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("./duplicates <dir> <log file>")
	}

	knownHashes = make(map[string]string)

	log, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Println("Scanning directory")
	recursivelyScanDirectory("", "/", os.Args[1], knownHashes, log)

	log.Close()
}
