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

var Debug = false

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
	outwriters := []io.Writer{stdoutbuf}
	errwriters := []io.Writer{stderrbuf}
	if Debug {
		outwriters = append(outwriters, os.Stdout)
		errwriters = append(outwriters, os.Stderr)
	}
	teeout := io.MultiWriter(outwriters...)
	teeerr := io.MultiWriter(errwriters...)

	go io.Copy(teeout, stdout)
	go io.Copy(teeerr, stderr)

	if err := cmd.Run(); err != nil {
		return Std{}, err
	}
	return Std{stdoutbuf.String(), stderrbuf.String()}, nil
}
