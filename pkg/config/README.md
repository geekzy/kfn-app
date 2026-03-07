# Config Package

The config package provides a centralized configuration system for the KFn FaaS platform using Viper for configuration management.

## Features

- YAML configuration file support
- Environment variable overrides
- Default values for all configuration options
- Hot reloading of configuration changes
- Structured configuration with validation

## Configuration Structure

The configuration is organized into logical sections:

- `registry`: Etcd-based function registry configuration
- `container`: Container runtime (containerd) configuration
- `sandbox`: Sandboxing (gVisor/Kata) configuration
- `image-registry`: Harbor registry configuration
- `observability`: Logging and metrics configuration
- `events`: Event system (Kafka/NATS) configuration
- `auth`: Authentication and authorization configuration
- `secrets`: Secrets management (Vault) configuration
- `build`: Image build system configuration
- `api-server`: HTTP API gateway configuration
- `scheduler`: Function scheduler configuration
- `runtime-wrappers`: Language runtime wrapper configuration

## Usage

```go
import "github.com/geekzy/kfn-app/pkg/config"

// Initialize configuration
if err := config.Init(); err != nil {
    log.Fatal(err)
}

// Get configuration instance
cfg := config.Get()

// Access configuration values
fmt.Println(cfg.Registry.Endpoints)
fmt.Println(cfg.Container.ContainerdSocket)
```

## Environment Variables

All configuration options can be overridden using environment variables with the `KFN_` prefix:

```bash
KFN_REGISTRY_ETCD_ENDPOINTS=localhost:2379,localhost:2380
KFN_CONTAINER_CONTAINERSOCKET=/run/containerd/containerd.sock
KFN_SANDBOX_TYPE=gVisor
```

## Configuration File Locations

The configuration system looks for files in the following locations:

1. Current directory (`./default.yaml`)
2. Configs directory (`./configs/default.yaml`)
3. System config directory (`/etc/kfn/default.yaml`)

## Hot Reloading

Configuration files are monitored for changes and automatically reloaded when modified.