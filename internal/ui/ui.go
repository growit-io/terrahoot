package ui

import "os"

type UserInterface interface {
	IsInteractive() bool
}

func New() UserInterface {
	if os.Getenv("CI") == "true" {
		return &ci{}
	} else {
		return &tty{}
	}
}
