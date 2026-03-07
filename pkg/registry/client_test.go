package registry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/geekzy/kfn-app/pkg/config"
)

func TestFunctionMetadata(t *testing.T) {
	// Test that FunctionMetadata struct is correctly defined
	fn := &FunctionMetadata{
		Name:           "test-function",
		Language:       "nodejs",
		RuntimeVersion: "18",
		RuntimeImage:   "node:18-alpine",
		Memory:         128,
		Handler:        "index.handler",
		State:          FunctionReady,
	}

	assert.Equal(t, "test-function", fn.Name)
	assert.Equal(t, "nodejs", fn.Language)
	assert.Equal(t, "18", fn.RuntimeVersion)
	assert.Equal(t, "node:18-alpine", fn.RuntimeImage)
	assert.Equal(t, 128, fn.Memory)
	assert.Equal(t, "index.handler", fn.Handler)
	assert.Equal(t, FunctionReady, fn.State)
}

func TestFunctionStateConstants(t *testing.T) {
	// Test that all function state constants are correctly defined
	assert.Equal(t, FunctionState("pending"), FunctionPending)
	assert.Equal(t, FunctionState("ready"), FunctionReady)
	assert.Equal(t, FunctionState("error"), FunctionError)
	assert.Equal(t, FunctionState("deleting"), FunctionDeleting)
}

func TestClient_FunctionKeyGeneration(t *testing.T) {
	// Test the helper function for generating function keys
	cfg := &config.Config{
		Registry: config.RegistryConfig{
			Prefix: "/kfn-test",
		},
	}

	client := &Client{
		config:    cfg,
		keyPrefix: "/kfn-test/",
	}

	expectedKey := "/kfn-test/functions/test-function"
	actualKey := client.functionKey("test-function")

	assert.Equal(t, expectedKey, actualKey)
}

func TestClient_JsonSerialization(t *testing.T) {
	// Test that FunctionMetadata can be serialized to JSON and back
	originalFn := &FunctionMetadata{
		Name:           "test-function",
		Language:       "nodejs",
		RuntimeVersion: "18",
		RuntimeImage:   "node:18-alpine",
		Memory:         128,
		Timeout:        30 * time.Second,
		Handler:        "index.handler",
		State:          FunctionReady,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Serialize to JSON
	data, err := json.Marshal(originalFn)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Deserialize from JSON
	var fn FunctionMetadata
	err = json.Unmarshal(data, &fn)
	assert.NoError(t, err)

	// Check that the values are the same
	assert.Equal(t, originalFn.Name, fn.Name)
	assert.Equal(t, originalFn.Language, fn.Language)
	assert.Equal(t, originalFn.RuntimeVersion, fn.RuntimeVersion)
	assert.Equal(t, originalFn.RuntimeImage, fn.RuntimeImage)
	assert.Equal(t, originalFn.Memory, fn.Memory)
	assert.Equal(t, originalFn.Handler, fn.Handler)
	assert.Equal(t, originalFn.State, fn.State)
}