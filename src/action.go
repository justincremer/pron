package pron

import (
	"errors"
	"os/exec"
	"strings"
)

// Represents a native go method
type internalAction struct {
	fn func(...interface{}) ([]byte, error)
}

// Represents an external shell command
type externalAction struct {
	cmd *exec.Cmd
}

// Common interface for dispatching actions
// type action interface {
// 	Dispatch() ([]byte, error)
// }

// // Internal action dispatch
// func (a *internalAction) Dispatch() ([]byte, error) {
// 	return a.fn()
// }

// // External action dispatch
// func (a *externalAction) Dispatch() ([]byte, error) {
// 	return a.cmd.Output()
// }

// Parses an external shell command
func parseExternalAction(i string) (a externalAction, err error) {
	switch argv := strings.Fields(i); len(argv) {
	case 0:
		return externalAction{}, errors.New("Command must not be nil")
	case 1:
		return externalAction{cmd: exec.Command(argv[0])}, nil
	default:
		return externalAction{cmd: exec.Command(argv[0], argv[1:]...)}, nil
	}
}
