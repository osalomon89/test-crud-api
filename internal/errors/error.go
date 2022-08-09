package itemerrors

import (
	"fmt"
)

type ItemError struct {
	Message string
}

func (e ItemError) Error() string {
	return fmt.Sprintf("item error: '%s'", e.Message)
}

type ResourceNotFoundError struct {
	Message string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}
