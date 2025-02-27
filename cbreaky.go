package termy

import (
	"errors"

	"golang.org/x/sys/unix"
)

type TermSettings struct {
	saved   unix.Termios
	isSaved bool
	fd      int
	echo    bool
}

// New creates an instance of TermSettings to handle the terminal state.
// fd: file descriptor (as in os.Stdout.Fd()),
// echo: Set to false to disable automatic echoing of user input.
func New(fd int, echo bool) *TermSettings {
	return &TermSettings{fd: fd, echo: echo}
}

// Cbreaky set the terminal (actually, STDIN) in a cbreak-like mode.
// Accepts an optional bool to disable echo as well.
// On success it saves the original state so you can retore it later.
// Example usage:
//
// term := termy.New(int(os.Stdin.Fd()), false)
// err := term.Cbreaky()
// if err != nil ...
// defer term.Restore()
func (ts *TermSettings) Cbreaky() error {
	termios, err := unix.IoctlGetTermios(ts.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	ts.saved = *termios
	ts.isSaved = true

	if !ts.echo {
		noEcho(termios)
	}
	noIcanon(termios)

	err = unix.IoctlSetTermios(ts.fd, unix.TIOCSETA, termios)
	if err != nil {
		return err
	}
	return nil
}

// Restore sets the terminal to its previous state.
// It returns an error if the previous state was not saved.
// Tipically you will call Restore after Cbreaky (probably with `defer`)
func (ts *TermSettings) Restore() error {
	if !ts.isSaved {
		return errors.New("err: terminal stated was not previously saved")
	}
	err := unix.IoctlSetTermios(ts.fd, unix.TIOCSETA, &ts.saved)
	if err != nil {
		return err
	}
	return nil
}

func noEcho(termios *unix.Termios) {
	termios.Lflag &^= unix.ECHO
}

func noIcanon(termios *unix.Termios) {
	termios.Lflag &^= unix.ICANON
}
