package repository

import (
	"io"
	"os/exec"
)

func Run(repo Repository, name string, args []string, stdout, stderr io.Writer) error {
	command := exec.Command(name, args...)
	command.Dir = repo.Directory
	command.Stdout = stdout
	command.Stderr = stderr

	return command.Run()
}
