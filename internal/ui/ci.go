package ui

type ci struct{}

func (ui *ci) IsInteractive() bool {
	return false
}
