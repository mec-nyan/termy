package term

import (
	"errors"

	"golang.org/x/sys/unix"
)

type msg struct {
	on bool
}

type Settings struct {
	// saved should only be written once and should be kept unchanged, so we can
	// restore the terminal to its previous settings at the end of our program.
	saved *unix.Termios
	// current is used to keep track of the terminal current state,
	// since we can turn some features on/off during the program execution.
	current unix.Termios
	fd      int
}

func New(fd int) *Settings {
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

	// Save the previous state of the terminal so we can restore it later.
	s.saved = termios
	// Save a COPY to keep track of current settings.
	s.current = *termios

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

func (s *Settings) Echoing() bool {
	flag := s.current.Lflag & unix.ECHO
	return flag == unix.ECHO
}

func (s *Settings) Cooked() bool {
	flag := s.current.Lflag & unix.ICANON
	return flag == unix.ICANON
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
// Typically you will call Restore after UnCookIt (probably with `defer`)
func (s *Settings) Restore() error {
	if s.saved == nil {
		return errors.New("err: terminal stated was not previously saved")
	}
	err := unix.IoctlSetTermios(s.fd, unix.TIOCSETA, s.saved)
	if err != nil {
		return err
	}
	return nil
}

// -------- Internal -------- //

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
	// Make a copy in case the operation fails.
	termios := s.current

	if echo != nil {
		if echo.on {
			setEcho(&termios)
		} else {
			setNoEcho(&termios)
		}
	}

	if icanon != nil {
		if icanon.on {
			setIcanon(&termios)
		} else {
			setNoIcanon(&termios)
		}
	}

	err := unix.IoctlSetTermios(s.fd, unix.TIOCSETA, &termios)
	if err != nil {
		return err
	}

	// Save the new state.
	s.current = termios

	return nil
}
