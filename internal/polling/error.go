package polling

import "fmt"

type ErrFetch struct {
	msg string
	err error
}

func (e *ErrFetch) Error() string {
	return fmt.Sprintf("fetch failed: %v. info: %s", e.err, e.msg)
}

type ErrParseData struct {
	data string
	err  error
}

func (e *ErrParseData) Error() string {
	return fmt.Sprintf("failed to parse document: %v: [%s]", e.err, e.data)
}
