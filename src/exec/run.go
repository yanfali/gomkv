package exec

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Std struct {
	Out string
	Err string
}

func Command(file string, params ...string) (Std, error) {
	name, err := exec.LookPath(file)
	if err != nil {
		return Std{}, err
	}
	stdoutbuf := bytes.NewBuffer([]byte{})
	stderrbuf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(name, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Std{}, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Std{}, err
	}
	teeout := io.MultiWriter(os.Stdout, stdoutbuf)
	teeerr := io.MultiWriter(os.Stderr, stderrbuf)

	go io.Copy(teeout, stdout)
	go io.Copy(teeerr, stderr)

	if err := cmd.Run(); err != nil {
		return Std{}, err
	}
	return Std{stdoutbuf.String(), stderrbuf.String()}, nil
}
