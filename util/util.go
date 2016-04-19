/*
Package util provides utility functions used by other packages.
*/
package util

import (
	"errors"
	"os/exec"
	"strings"
)

// ErrDetached is returned if it looks like a Git repo is in a
// detached state.
var ErrDetached = errors.New("Your repo looks like it is in a detached state")

// GitBranch returns the current git branch for the given directory
func GitBranch(dir string) (branch string, err error) {
	cmd := Cmd(dir, "git rev-parse --abbrev-ref HEAD")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	branch = strings.TrimSpace(string(out))
	if branch == "HEAD" {
		err = ErrDetached
	}
	return
}

// Cmd creates an exec.Cmd to run in the given directory.
func Cmd(dir string, command string) *exec.Cmd {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	return cmd
}

// CmdOutput represents the output of a bash command.
type CmdOutput struct {
	Cmd    string
	Output string
	Err    error
}

// RunAll runs all the given commands in the given directory, and
// returns a list of CmdOutputs.
func RunAll(dir string, commands ...string) []*CmdOutput {
	outputs := make([]*CmdOutput, len(commands))
	var i int
	var command string
	for i, command = range commands {
		out, err := Cmd(dir, command).Output()
		outputs[i] = &CmdOutput{command, string(out), err}
		if err != nil {
			break
		}
	}
	return outputs[0 : i+1]
}
