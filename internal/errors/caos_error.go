package errors

import (
	"fmt"
)

var _ Error = (*CaosError)(nil)

type CaosError struct {
	Parent  error
	Message string
	ID      string
}

func ThrowError(parent error, id, message string) error {
	return CreateCaosError(parent, id, message)
}

func CreateCaosError(parent error, id, message string) *CaosError {
	return &CaosError{
		Parent:  parent,
		ID:      id,
		Message: message,
	}
}

func (err *CaosError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("AggregateID=%s Message=%s Parent=(%v)", err.ID, err.Message, err.Parent)
	}
	return fmt.Sprintf("AggregateID=%s Message=%s", err.ID, err.Message)
}

func (err *CaosError) Unwrap() error {
	return err.GetParent()
}

func (err *CaosError) GetParent() error {
	return err.Parent
}

func (err *CaosError) GetMessage() string {
	return err.Message
}

func (err *CaosError) GetID() string {
	return err.ID
}
