package papersave

import (
	"os/exec"
	"os"
	"io"
	"fmt"
	"bytes"
	"sync"
)

// From: https://github.com/kjk/go-cookbook/blob/master/advanced-exec/03-live-progress-and-capture-v3.go
func runCommand(dir string, cmdLine ...string) {
	pwd, err := os.Getwd()
	defer os.Chdir(pwd)
	err = os.Chdir(dir)
	Panicp(err)
	cmd := 	exec.Command(cmdLine[0], cmdLine[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Start()
	Panicp(err)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	Panicp(err)
	if errStdout != nil || errStderr != nil {
		Fatal("failed to capture stdout or stderr\n")
	}
	if false {
		outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
		
		fmt.Println(cmd)
	}
}
