package errors

import "fmt"

type EntityNotFound struct {
	Entity string
	ID     int
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("entity %s with id %d not found", e.Entity, e.ID)
}
