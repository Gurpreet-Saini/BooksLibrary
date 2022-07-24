package errors

import (
	"fmt"
)

type InValidDetails struct {
	Details string
}

func (e InValidDetails) Error() string {
	return fmt.Sprintf("detail %s is invalid", e.Details)
}
