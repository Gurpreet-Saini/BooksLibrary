package errors

import "fmt"

type ExistAlready struct {
	Entity string
}

func (e ExistAlready) Error() string {
	return fmt.Sprintf("entity  %v already exists", e.Entity)
}
