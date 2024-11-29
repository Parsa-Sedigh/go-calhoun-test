package git

import (
	"os/exec"
	"strings"
)

var execCommand = exec.Command

// Version returns the current git version that the user is using
func Version() string {
	cmd := execCommand("git", "version")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	n := len("git version ")
	version := string(stdout[n:])

	return strings.TrimSpace(version)
}

type Checker struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

// command sets the execCommand field if it's not set, otherwise just return it if it was injected into Checker.
func (gc *Checker) command(name string, arg ...string) *exec.Cmd {
	if gc.execCommand == nil {
		return exec.Command(name, arg...)
	}

	return gc.execCommand(name, arg...)
}

func (gc *Checker) Version() string {
	cmd := gc.command("git", "version")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	n := len("git version ")
	version := string(stdout[n:])
	return strings.TrimSpace(version)
}
