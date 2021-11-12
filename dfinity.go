package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

func deployWithDfx(targetpath string, f *os.File) error {

	deploycmd := exec.Command("dfx", "deploy", "--network", "ic")
	deploycmd.Dir = targetpath

	stderr, err := deploycmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := deploycmd.StdoutPipe()
	if err != nil {
		return err
	}

	deploycmd.Start()

	errReader := bufio.NewReader(stderr)
	outReader := bufio.NewReader(stdout)

	for {
		line, err := errReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// write local
		_, err = f.WriteString(Format(line))
		if err != nil {
			break
		}
	}

	for {
		line, err := outReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// write local
		_, err = f.WriteString(Format(line))
		if err != nil {
			break
		}
	}

	deploycmd.Wait()

	return nil
}
