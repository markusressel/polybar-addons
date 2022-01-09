package util

import (
	"bytes"
	"os"
	"os/exec"
)

// ExecCommand Executes a shell command with the given arguments
// and returns its stdout as a []byte.
// If an error occurs the content of stderr is printed
// and an error is returned.
func ExecCommand(command string, args ...string) (string, error) {
	//log.Printf("Executing command: %s %s", command, args)
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return string(stdout.Bytes()), nil
}

// ExecCommandEnv Like execCommand but with the possibility to add environment variables
// to the executed process.
func ExecCommandEnv(env []string, attach bool, command string, args ...string) (string, error) {
	//log.Printf("Executing command: %s %s", command, args)
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	var err error
	if attach {
		err = cmd.Run()
	} else {
		err = cmd.Start()
		if err != nil {
			return "", err
		}
		err = cmd.Process.Release()
	}

	if err != nil {
		return "", err
	}

	return string(stdout.Bytes()), nil
}
