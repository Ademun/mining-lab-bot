package errs

import (
	"errors"
	"fmt"
)

var ErrSubscriptionExists = errors.New("subscription already exists")

type ErrQueryExecution struct {
	Operation string
	Query     string
	Err       error
}

func (e *ErrQueryExecution) Error() string {
	return fmt.Sprintf("failed to execute Query during: %s: [%s] %v", e.Operation, e.Query, e.Err)
}

type ErrRowIteration struct {
	Operation string
	Query     string
	Err       error
}

func (e *ErrRowIteration) Error() string {
	return fmt.Sprintf("failed to iterate rows during: %s: [%s] %v", e.Operation, e.Query, e.Err)
}
