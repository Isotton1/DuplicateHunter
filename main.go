package main

import (
	"crypto/sha512"
	"fmt"
	"io/fs"
	"log"
	"os"
)

//TODO Organize this mess.

type fileInfo struct {
	path string
	hash [64]byte
	size int
}

var fileStack []fileInfo

const root_dir = "/home/linuxsupremacy/coding/NoMoreDup/"

func main() {
	fmt.Println(root_dir)
	fs.WalkDir(os.DirFS(root_dir), ".", newFileInfo)
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
	fmt.Printf("start walk in %s %s \n", root_dir, path)
	//if is a dir
	if err != nil {
		return fs.SkipDir
	}
	if dir.IsDir() {
		fmt.Println("a dir")
		if dir.Name() == ".." {
			fmt.Println("skipped")
			return fs.SkipDir
		}
		if dir.Name()[0] == '.' {
			fmt.Println(". dir")
			return nil
		}
		return fs.WalkDir(os.DirFS(root_dir), dir.Name(), newFileInfo)
	}
	//skip dot file.
	if dir.Name()[0] == '.' {
		return nil
	}
	//Open file
	file, err := fs.ReadFile(os.DirFS(root_dir), dir.Name())
	if err != nil {
		log.Panic(err)
	}
	//Create hash
	hash := sha512.Sum512(file)
	//Calculate size
	size := len(file)
	//Add fileInfo to fileStack
	curFileInfo := fileInfo{
		path: root_dir + path,
		hash: hash,
		size: size,
	}
	fmt.Println(curFileInfo.path)
	fileStack = append(fileStack, curFileInfo)
	return nil
}
