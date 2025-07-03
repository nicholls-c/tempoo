package internal

import "fmt"

type TempooError struct {
	Message string
	Cause   error
}

// error returns the error message
func (e *TempooError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// InvalidIssueKeyError is an error type for invalid issue keys
type InvalidIssueKeyError struct {
	IssueKey string
}

// error returns the error message
func (e *InvalidIssueKeyError) Error() string {
	return fmt.Sprintf("Issue key %s is not valid", e.IssueKey)
}
