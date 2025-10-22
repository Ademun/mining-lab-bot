package polling

import "fmt"

type ErrFetch struct {
	url string
	msg string
	err error
}

func (e *ErrFetch) Error() string {
	return fmt.Sprintf("fetch failed for %s: %v. info: %s", e.url, e.err, e.msg)
}

type ErrParseData struct {
	data string
	msg  string
	err  error
}

func (e *ErrParseData) Error() string {
	return fmt.Sprintf("failed to parse document: %v: [%s]. info: %s", e.err, e.data, e.msg)
}
