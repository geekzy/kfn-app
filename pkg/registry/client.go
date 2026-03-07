package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/geekzy/kfn-app/pkg/config"
)

// Client is the etcd registry client
type Client struct {
	client    *clientv3.Client
	config    *config.Config
	keyPrefix string
}

// NewClient creates a new etcd registry client
func NewClient(cfg *config.Config) (*Client, error) {
	// Create etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Registry.Endpoints,
		DialTimeout: cfg.Registry.DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	// Ensure prefix ends with "/"
	prefix := cfg.Registry.Prefix
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	return &Client{
		client:    cli,
		config:    cfg,
		keyPrefix: prefix,
	}, nil
}

// Close closes the etcd client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// CreateFunction creates a new function in the registry
func (c *Client) CreateFunction(ctx context.Context, fn *FunctionMetadata) error {
	// Set timestamps
	now := time.Now()
	fn.CreatedAt = now
	fn.UpdatedAt = now

	// Serialize function metadata
	data, err := json.Marshal(fn)
	if err != nil {
		return fmt.Errorf("failed to marshal function metadata: %w", err)
	}

	// Put function metadata in etcd
	key := c.functionKey(fn.Name)
	_, err = c.client.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("failed to put function in etcd: %w", err)
	}

	return nil
}

// GetFunction retrieves a function from the registry
func (c *Client) GetFunction(ctx context.Context, name string) (*FunctionMetadata, error) {
	// Get function metadata from etcd
	key := c.functionKey(name)
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get function from etcd: %w", err)
	}

	// Check if function exists
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// Deserialize function metadata
	var fn FunctionMetadata
	if err := json.Unmarshal(resp.Kvs[0].Value, &fn); err != nil {
		return nil, fmt.Errorf("failed to unmarshal function metadata: %w", err)
	}

	return &fn, nil
}

// ListFunctions lists all functions in the registry
func (c *Client) ListFunctions(ctx context.Context) ([]*FunctionMetadata, error) {
	// Get all functions from etcd
	prefix := c.functionsPrefix()
	resp, err := c.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to list functions from etcd: %w", err)
	}

	// Deserialize function metadata
	functions := make([]*FunctionMetadata, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var fn FunctionMetadata
		if err := json.Unmarshal(kv.Value, &fn); err != nil {
			// Skip invalid entries
			continue
		}
		functions = append(functions, &fn)
	}

	return functions, nil
}

// UpdateFunction updates a function in the registry
func (c *Client) UpdateFunction(ctx context.Context, fn *FunctionMetadata) error {
	// Update timestamp
	fn.UpdatedAt = time.Now()

	// Serialize function metadata
	data, err := json.Marshal(fn)
	if err != nil {
		return fmt.Errorf("failed to marshal function metadata: %w", err)
	}

	// Put function metadata in etcd
	key := c.functionKey(fn.Name)
	_, err = c.client.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("failed to update function in etcd: %w", err)
	}

	return nil
}

// DeleteFunction deletes a function from the registry
func (c *Client) DeleteFunction(ctx context.Context, name string) error {
	// Delete function from etcd
	key := c.functionKey(name)
	_, err := c.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete function from etcd: %w", err)
	}

	return nil
}

// RegisterFunctionWithLease registers a function with a lease for health tracking
func (c *Client) RegisterFunctionWithLease(ctx context.Context, name string, ttl int64) error {
	// Create a lease
	leaseResp, err := c.client.Grant(ctx, ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	// Create a key for health tracking
	key := c.healthKey(name)
	value := "alive"

	// Put key with lease
	_, err = c.client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to register function with lease: %w", err)
	}

	return nil
}

// KeepAliveFunction extends the lease for a function
func (c *Client) KeepAliveFunction(ctx context.Context, name string, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	// Keep alive the lease
	ch, err := c.client.KeepAlive(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to keep alive function: %w", err)
	}

	return ch, nil
}

// AcquireLock acquires a distributed lock for a function
func (c *Client) AcquireLock(ctx context.Context, name string) (*concurrency.Mutex, error) {
	// Create a session
	session, err := concurrency.NewSession(c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create a mutex
	mutex := concurrency.NewMutex(session, c.lockKey(name))

	// Lock
	if err := mutex.Lock(ctx); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return mutex, nil
}

// helper functions

// functionKey returns the etcd key for a function
func (c *Client) functionKey(name string) string {
	return path.Join(c.keyPrefix, "functions", name)
}

// functionsPrefix returns the etcd prefix for all functions
func (c *Client) functionsPrefix() string {
	return path.Join(c.keyPrefix, "functions") + "/"
}

// healthKey returns the etcd key for function health tracking
func (c *Client) healthKey(name string) string {
	return path.Join(c.keyPrefix, "health", name)
}

// lockKey returns the etcd key for function locks
func (c *Client) lockKey(name string) string {
	return path.Join(c.keyPrefix, "locks", name)
}