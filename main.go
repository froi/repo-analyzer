package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// This may not be needed.
func createFile(fileName string) (*os.File, error) {
	createdFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] There was an error creating", fileName, err)
		return nil, err
	}

	defer func() {
		err := createdFile.Close()
		if err != nil {
			fmt.Println("[ERROR] There was an error closing the file: ", fileName)
			return nil, err
		}
	}()

	return createdFile, nil
}

// Execute git lfs track on the supplied extension
func addFileToLfs(fileExtension string) error {
	cmd := exec.Command("git", "lfs", "track", fmt.Sprintf("*.%s", fileExtension))
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}
func main() {

	gitAttr, err := createFile(".gitattributes")
	if err != nil {
		log.Fatal("[ERROR] Creating .gitattributes file", err)
	}

	var myWalkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fmt.Println("current path;", path)
			fmt.Println("FileInfo.Name:", info.Name())
			fmt.Println("Size:", info.Size())

			if info.Size()/(1024*1000) > 40 {
				// Add file ext to git lfs
				ext := filepath.Ext(info.Name())
				err := addFileToLfs(ext)
				if err != nil {
					return err
				}
			}
			cmd := exec.Command("git", "add", path)
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := filepath.Walk(".", myWalkFunc)
	if err != nil {
		log.Fatal("Error walking the path", err)
	}
}
