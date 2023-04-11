package user

import "fmt"

type UserError struct {
	Message string
}

func (e UserError) Error() string {
	return fmt.Sprintf("user error: '%s'", e.Message)
}

type ResourceNotFoundError struct {
	Message string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}
