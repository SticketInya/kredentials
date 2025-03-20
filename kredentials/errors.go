package kredentials

import (
	"fmt"
)

type ErrKredentialConflict struct {
	Name string
}

func (e ErrKredentialConflict) Error() string {
	return fmt.Sprintf("kredential '%s' already exists", e.Name)
}
