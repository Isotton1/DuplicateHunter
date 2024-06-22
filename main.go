package main

import (
	"crypto/sha512"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// TODO
//   - Write -h.
//   - Implement delete.
type fileInfo struct {
	path string
	hash [64]byte
	size int
}

var fileStack []fileInfo

var config struct {
	delete bool
	oldest bool
}

func main() {
	argv := os.Args
	argc := len(argv)

	if argc > 2 {
		for i := 1; i < argc; i++ {
			switch argv[i] {
			case "-h":
				fmt.Println("")
			case "-d":
				config.delete = true
			case "-do":
				config.delete = true
				config.oldest = true
			default:
				log.Fatal(argv[i] + " is not a valid flag")
			}
		}
	}

	filepath.WalkDir(".", newFileInfo)
	compareFiles()
	fmt.Println("finish")
}

func compareFiles() {
	for i := 0; i < len(fileStack)-1; i++ {
		for j := i + 1; j < len(fileStack); j++ {
			if fileStack[i].size == fileStack[j].size &&
				fileStack[i].hash == fileStack[j].hash {
				fmt.Println(fileStack[i].path + " is equal to " + fileStack[j].path)
			}
		}
	}
}

func newFileInfo(path string, dir fs.DirEntry, err error) error {
	fmt.Printf("start walk in %s%s \n", ".", path)
	if err != nil {
		return fs.SkipDir
	}
	//skip dot file/dir and don't skip "."
	if dir.Name()[0] == '.' && len(dir.Name()) > 1 {
		return fs.SkipDir
	}
	if dir.IsDir() {
		return nil
	}
	curFileInfo, err := createFileInfo(path)
	if err != nil {
		log.Panic(err)
	}
	fileStack = append(fileStack, curFileInfo)
	return nil
}

func createFileInfo(path string) (fileInfo, error) {
	//Open file
	file, err := fs.ReadFile(os.DirFS("."), path)
	if err != nil {
		return fileInfo{}, err
	}
	//Create hash
	hash := sha512.Sum512(file)
	//Calculate size
	size := len(file)
	//Add fileInfo to fileStack
	curFileInfo := fileInfo{
		path: filepath.Join(".", path),
		hash: hash,
		size: size,
	}
	return curFileInfo, nil
}
