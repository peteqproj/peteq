package automation

import "github.com/peteqproj/peteq/pkg/tenant"

type (
	// Automation that start some logical workflow
	Automation struct {
		tenant.Tenant `json:"tenant" yaml:"tenant"`
		Metadata      Metadata `json:"metadata" yaml:"metadata"`
		Spec          Spec     `json:"spec" yaml:"spec"`
	}

	// Metadata of automation
	Metadata struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
	}

	// Spec for automation
	// can be have on of types crontab or url
	Spec struct {
		Type            string `json:"type"` // automation type task.archiver
		JSONInputSchema string `json:"jsonInputSchema"`
	}

	// TriggerBinding trigger->automation
	TriggerBinding struct {
		tenant.Tenant `json:"tenant" yaml:"tenant"`
		Metadata      TriggerBindingMetadata `json:"metadata" yaml:"metadata"`
		Spec          TriggerBindingSpec     `json:"spec" yaml:"spec"`
	}

	TriggerBindingMetadata struct {
		ID   string `json:"id" yaml:"id"`
		Name string `json:"name" yaml:"name"`
	}

	TriggerBindingSpec struct {
		Automation string `json:"automation" yaml:"automation"`
		Trigger    string `json:"trigger" yaml:"trigger"`
	}
)
