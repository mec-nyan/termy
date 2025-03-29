package termy

import (
	"errors"

	"golang.org/x/sys/unix"
)

type msg struct {
	on bool
}

type Settings struct {
	saved   unix.Termios
	isSaved bool
	fd      int
}

func NewTerm(fd int) *Settings {
	return &Settings{
		fd: fd,
	}
}

func (s *Settings) Init() error {
	// Get the current state of the terminal.
	// Actually, this will get the configuration for the file descriptor.
	termios, err := unix.IoctlGetTermios(s.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}

	s.saved = *termios
	s.isSaved = true

	return nil
}

func (s *Settings) CookIt() error {
	icanon := msg{true}
	return s.set(&icanon, nil)
}

func (s *Settings) UnCookIt() error {
	icanon := msg{false}
	return s.set(&icanon, nil)
}

func (s *Settings) Echo() error {
	echo := msg{true}
	return s.set(nil, &echo)
}

func (s *Settings) NoEcho() error {
	echo := msg{false}
	return s.set(nil, &echo)
}

// TODO: How should we handle "resize"?
func (s *Settings) Size() (rows, cols int, err error) {
	ws, err := unix.IoctlGetWinsize(s.fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Row), int(ws.Col), nil
}

// Restore sets the terminal to its previous state.
// It returns an error if the previous state was not saved.
// Tipically you will call Restore after Cbreaky (probably with `defer`)
func (s *Settings) Restore() error {
	if !s.isSaved {
		return errors.New("err: terminal stated was not previously saved")
	}
	err := unix.IoctlSetTermios(s.fd, unix.TIOCSETA, &s.saved)
	if err != nil {
		return err
	}
	return nil
}

func setEcho(termios *unix.Termios) {
	termios.Lflag |= unix.ECHO
}

func setNoEcho(termios *unix.Termios) {
	termios.Lflag &^= unix.ECHO
}

func setIcanon(termios *unix.Termios) {
	termios.Lflag |= unix.ICANON
}

func setNoIcanon(termios *unix.Termios) {
	termios.Lflag &^= unix.ICANON
}

func (s *Settings) set(icanon, echo *msg) error {
	// Get the current state of the terminal.
	// Actually, this will get the configuration for the file descriptor.
	termios, err := unix.IoctlGetTermios(s.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}

	if echo != nil {
		if echo.on {
			setEcho(termios)
		} else {
			setNoEcho(termios)
		}
	}

	if icanon != nil {
		if icanon.on {
			setIcanon(termios)
		} else {
			setNoIcanon(termios)
		}
	}

	err = unix.IoctlSetTermios(s.fd, unix.TIOCSETA, termios)
	if err != nil {
		return err
	}

	return nil
}
