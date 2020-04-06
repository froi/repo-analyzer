package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Git struct {
	branch string
	cmd    string
	err    error
}

func executeCommand(commandString string, subcommand string, cmdArgs []string) error {

	cmdArgs = append([]string{subcommand}, cmdArgs...)
	cmd := exec.Command(commandString, cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out

	return cmd.Run()
}

func (git *Git) Init(commandArgs ...string) *Git {
	if git.err != nil {
		return git
	}

	git.err = executeCommand(git.cmd, "init", commandArgs)
	return git
}
func (git *Git) Push(commandArgs ...string) *Git {
	if git.err != nil {
		return git
	}
	git.err = executeCommand(git.cmd, "push", commandArgs)
	return git
}

func (git *Git) Add(commandArgs ...string) *Git {
	if git.err != nil {
		return git
	}
	git.err = executeCommand(git.cmd, "add", commandArgs)

	return git
}

func (git *Git) Commit(commandArgs ...string) *Git {
	if git.err != nil {
		return git
	}

	git.err = executeCommand(git.cmd, "commit", commandArgs)

	return git
}

func (git *Git) AddToLfs(fileExtension string) *Git {
	if git.err != nil {
		return git
	}
	cmdArgs := []string{"track", fmt.Sprintf("*.%s", fileExtension)}
	git.err = executeCommand(git.cmd, "lfs", cmdArgs)

	return git
}

func (git *Git) Error() error {
	return git.err
}
