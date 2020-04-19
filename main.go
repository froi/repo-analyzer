package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

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
				Push().
				Error()
			if err != nil {
				return err
			}
		}
		count++
	}
	return git.Commit("-m", "Adding what's left of the non LFS files").
		Push().
		Error()
}

func configureGitLfs(git *Git, files *[]File) {
	for _, file := range *files {

		err := git.AddToLfs(file.ext).
			Add().
			Commit("-m", fmt.Sprintf("Add extension %s to Git LFS tracking", file.ext)).
			Push().
			Error()
		if err != nil {
			log.Println("There was a Git LFS error: tracking an extension", file.ext)
		}
	}
}
func processLfsFiles(git *Git, files *[]File, branch string, remoteName string) error {

	for _, file := range *files {
		err := git.Add(file.path).
			Commit("-m", "Adding LFS file to repo").
			Push().
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
func main() {
	repoInit := flag.Bool("init", false, "Specifies with a Git repo needs to be created.")
	remoteName := flag.String("remote", "origin", "Specifies the Git remote to use when pushing changes. Defaults to origin.")
	branch := flag.String("branch", "master", "Specifies the Git branch to use when commiting changes. Defaults to master.")
	lfsInstall := flag.Bool("lfs-instal", false, "Ensures that Git LFS hooks are installed for this repository.")
	flag.Parse()

	lfsFiles := make([]File, 0, 0)
	nonLfsFiles := make([]File, 0, 0)

	git := Git{
		branch:     *branch,
		cmd:        "git",
		err:        nil,
		remoteName: *remoteName,
	}

	if *repoInit {
		err := git.Init().Error()
		if err != nil {
			log.Fatal("Error initializing the Git repository", err)
		}
	} else {
		verifyGitRepo()
	}

	if *lfsInstall {
		err := git.LfsInstall().Error()
		if err != nil {
			log.Fatal("Error installing Git LFS hooks.", err)
		}
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
	configureGitLfs(&git, &lfsFiles)
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
