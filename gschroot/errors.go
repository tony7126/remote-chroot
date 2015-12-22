package gschroot

import (
	"fmt"
)
// signal error (they have their own error codes)
// http error (create an error code)

type TarNotFound struct {
	Url string
	Code int
}

func (t *TarNotFound) Error() string {
	return fmt.Sprintf("No Tar Found at URL: %s", t.Url)
}

type HTTPProtocolError struct {
	Msg string
	Code int
}

func (h *HTTPProtocolError) Error() string {
	return fmt.Sprintf("%s", h.Msg)
}

type InvalidCommandError struct {
	Msg string
	Code int
}

func (h *InvalidCommandError) Error() string {
	return fmt.Sprintf("%s", h.Msg)
}
