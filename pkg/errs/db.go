package errs

import (
	"errors"
	"fmt"
)

var ErrSubscriptionExists = errors.New("subscription already exists")

var ErrBeginTransaction = errors.New("failed to begin transaction")

type ErrQueryCreation struct {
	Operation string
	Query     string
	Err       error
}

func (e *ErrQueryCreation) Error() string {
	return fmt.Sprintf("failed to create query during: %s: [%s] %v", e.Operation, e.Query, e.Err)
}

type ErrQueryExecution struct {
	Operation string
	Query     string
	Err       error
}

func (e *ErrQueryExecution) Error() string {
	return fmt.Sprintf("failed to execute query during: %s: [%s] %v", e.Operation, e.Query, e.Err)
}

type ErrRowIteration struct {
	Operation string
	Query     string
	Err       error
}

func (e *ErrRowIteration) Error() string {
	return fmt.Sprintf("failed to iterate rows during: %s: [%s] %v", e.Operation, e.Query, e.Err)
}
