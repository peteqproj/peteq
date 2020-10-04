package user

type (
	// User of the system
	User struct {
		Metadata Metadata `json:"metadata" yaml:"metadata"`
		Spec     Spec     `json:"spec" yaml:"spec"`
	}

	// Metadata of user
	Metadata struct {
		ID    string `json:"id" yaml:"id"`
		Email string `json:"email" yaml:"email"`
	}

	// Spec of user
	Spec struct {
		TokenHash    string `json:"tokenHash" yaml:"tokenHash"`
		PasswordHash string `json:"passwordHash" yaml:"passwordHash"`
	}
)
