package tty

import "os"

// TTY provides the means to interact with stdout, stdin and stderr
type TTY struct {
	Stdout *os.File
	Stdin  *os.File
	Stderr *os.File
}

// New creates a new TTY with default values for stdin, stdout and stderr.
func New() TTY {
	return TTY{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
}
