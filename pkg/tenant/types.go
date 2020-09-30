package tenant

const (
	// User request key
	User Type = iota
)

// Type is used to set the key of context in request
type Type int

func (r Type) String() string {
	return [...]string{"User"}[r]
}
