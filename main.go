package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// This may not be needed.
func createFile(fileName string) (*os.File, error) {
	createdFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] There was an error creating", fileName, err)
		return nil, err
	}

	defer func() error {
		err := createdFile.Close()
		if err != nil {
			fmt.Println("[ERROR] There was an error closing the file: ", fileName)
			return err
		}
		return nil
	}()

	return createdFile, nil
}

type File struct {
	path string
	size int64
	name string
	ext  string
}

func processNonLfsFiles(git *Git, files *[]File, branch string, remoteName string) error {
	count := 0

	for _, file := range *files {
		git.Add(file.path)
		if count%10 == 0 {
			err := git.Commit("-m", "Adding non LFS files to repo").
				Push(remoteName, branch).
				Error()
			if err != nil {
				return err
			}
		}
		count++
	}
	return git.Commit("-m", "Adding what's left of the non LFS files").
		Push(branch).
		Error()
}

func processLfsFiles(git *Git, files *[]File, branch string, remoteName string) error {

	for _, file := range *files {
		err := git.Add(file.path).
			Commit("-m", "Adding LFS file to repo").
			Push(remoteName, branch).
			Error()
		if err != nil {
			return err
		}
	}
	return nil
}
func verifyGitRepo() {
	if _, err := os.Stat(".git"); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Not in a Git repository. Please us in a repository or initialize a new one.")
		} else {
			log.Fatal("There was an error verifying Git repo existence.", err)
		}
	}
}
func init() {

}
func main() {
	repoInit := flag.Bool("init", false, "Specifies with a Git repo needs to be created.")
	remoteName := flag.String("remote", "origin", "Specifies the Git remote to use when pushing changes. Defaults to origin.")
	branch := flag.String("branch", "master", "Specifies the Git branch to use when commiting changes. Defaults to master.")
	flag.Parse()

	lfsFiles := make([]File, 0, 0)
	nonLfsFiles := make([]File, 0, 0)

	git := Git{
		branch: "",
		cmd:    "git",
		err:    nil,
	}

	if *repoInit {
		err := git.Init().Error()
		if err != nil {
			log.Fatal("Error initializing the Git repository", err)
		}
	} else {
		verifyGitRepo()
	}

	var myWalkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file := File{
				name: info.Name(),
				path: path,
				size: info.Size(),
				ext:  filepath.Ext(info.Name()),
			}

			if info.Size()/(1024*1000) > 40 {
				lfsFiles = append(lfsFiles, file)
			} else {
				nonLfsFiles = append(nonLfsFiles, file)
			}
		}
		return nil
	}
	err := filepath.Walk(".", myWalkFunc)
	if err != nil {
		log.Fatal("Error walking the path", err)
	}
	// Process LFS
	err = processLfsFiles(&git, &lfsFiles, *branch, *remoteName)
	if err != nil {
		log.Fatal("Error processing LFS files", err)
	}
	// Process non LFS
	err = processNonLfsFiles(&git, &nonLfsFiles, *branch, *remoteName)
	if err != nil {
		log.Fatal("Error processing non LFS files", err)
	}
}
