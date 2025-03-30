package userservice

import (
	"fmt"
)

const (
	ServerError = 500
)

type RichErr struct {
	Message string
	Kind    int
}

func (err *RichErr) Error() string {
	switch err.Kind {
	case ServerError:
		return fmt.Sprintf("Server Error: %s", err.Message)
	default:
		return fmt.Sprintf("Not impelemented: %s", err.Message)
	}
}
