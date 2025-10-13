package subscription

import (
	"errors"
	"fmt"
)

var ErrSubscriptionExists = errors.New("subscription already exists")

type ErrQueryExecution struct {
	operation string
	query     string
	err       error
}

func (e *ErrQueryExecution) Error() string {
	return fmt.Sprintf("failed to execute query during: %s: [%s] %v", e.operation, e.query, e.err)
}

type ErrRowIteration struct {
	operation string
	query     string
	err       error
}

func (e *ErrRowIteration) Error() string {
	return fmt.Sprintf("failed to iterate rows during: %s: [%s] %v", e.operation, e.query, e.err)
}
