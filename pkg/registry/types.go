package registry

import (
	"time"
)

// FunctionState represents the state of a function
type FunctionState string

const (
	FunctionPending  FunctionState = "pending"
	FunctionReady    FunctionState = "ready"
	FunctionError    FunctionState = "error"
	FunctionDeleting FunctionState = "deleting"
)

// FunctionMetadata represents the metadata for a function
type FunctionMetadata struct {
	// Name is the unique identifier for the function
	Name string `json:"name" mapstructure:"name"`

	// Language is the programming language of the function (e.g., "nodejs", "python")
	Language string `json:"language" mapstructure:"language"`

	// RuntimeVersion is the version of the runtime (e.g., "18" for Node.js 18)
	RuntimeVersion string `json:"runtimeVersion" mapstructure:"runtimeVersion"`

	// RuntimeImage is the Docker image to use for the function
	RuntimeImage string `json:"runtimeImage" mapstructure:"runtimeImage"`

	// Memory is the memory limit for the function in MB
	Memory int `json:"memory" mapstructure:"memory"`

	// Timeout is the maximum execution time for the function
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`

	// Handler is the entry point for the function (e.g., "index.handler" for Node.js)
	Handler string `json:"handler" mapstructure:"handler"`

	// State is the current state of the function
	State FunctionState `json:"state" mapstructure:"state"`

	// CreatedAt is the timestamp when the function was created
	CreatedAt time.Time `json:"createdAt" mapstructure:"createdAt"`

	// UpdatedAt is the timestamp when the function was last updated
	UpdatedAt time.Time `json:"updatedAt" mapstructure:"updatedAt"`
}
