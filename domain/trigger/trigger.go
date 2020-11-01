package trigger

import "github.com/peteqproj/peteq/pkg/tenant"

type (
	// Trigger that start some logical workflow
	Trigger struct {
		tenant.Tenant `json:"tenant" yaml:"tenant"`
		Metadata      Metadata `json:"metadata" yaml:"metadata"`
		Spec          Spec     `json:"spec" yaml:"spec"`
	}

	// Metadata of project
	Metadata struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
	}

	// Spec for trigger
	// can be have on of types crontab or url
	Spec struct {
		Cron    *string  `json:"cron,omitempty"`
		Webhook *Webhook `json:"webhook,omitempty"`
	}

	// Webhook for trigger
	Webhook struct {
		URL             string
		RequiredHeaders map[string]string
	}
)
