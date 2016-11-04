package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// Std out and error output buffers
type Std struct {
	Out string
	Err string
}

// Debug global flag
var Debug = false

// Command with default timeout
func Command(file string, params ...string) (Std, error) {
	return CommandWithTimeout(file, time.Second*15, params...)
}

// CommandWithTimeout wrapper
func CommandWithTimeout(file string, timeout time.Duration, params ...string) (Std, error) {
	name, err := exec.LookPath(file)
	if err != nil {
		return Std{}, fmt.Errorf("Error can't find executable %q %v", file, err)
	}
	stdoutbuf := bytes.NewBuffer([]byte{})
	stderrbuf := bytes.NewBuffer([]byte{})
	ctxt, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctxt, name, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Std{}, fmt.Errorf("Error attaching to stdout: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Std{}, fmt.Errorf("Error attaching to stderr: %v", err)
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
		return Std{}, fmt.Errorf("Error running %q %v", file, err)
	}
	return Std{stdoutbuf.String(), stderrbuf.String()}, nil
}
