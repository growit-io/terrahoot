package ui

type tty struct{}

func (ui *tty) IsInteractive() bool {
	// XXX: depends on whether stdin, stdout or stderr is connected to a TTY
	return true
}
