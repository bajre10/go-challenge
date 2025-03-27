package appErrors

import "fmt"

type ErrorType int

const (
	Validation ErrorType = iota
)

type DBError struct {
	Field   string
	Type    ErrorType
	Message string
}

func (e *DBError) Error() string {
	return fmt.Sprintf("Error %d on field %s: %s", e.Type, e.Field, e.Message)
}
