package utils

import (
	"io"
	"os"
	"os/exec"
)

func SysRun(stdout io.Writer, stderr io.Writer, name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = stdout
	command.Stderr = stderr
	return command.Run()
}

func CdWithFunc(dir string, fu func()) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	fu()
	return nil
}

func TempCdWithFunc(dir string, fu func()) error {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	return CdWithFunc(dir, fu)
}
