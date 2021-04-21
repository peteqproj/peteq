package errors

import "fmt"

var ErrMissingUserInContext = fmt.Errorf("user missing in context")
