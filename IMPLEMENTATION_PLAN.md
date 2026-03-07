# KFn FaaS Platform - Implementation Plan

Created: 2026-03-01

## Overview

This document outlines the step-by-step implementation of the KFN Function-as-a-Service platform, a Go-based FaaS system supporting Node.js and Python functions with container-per-invocation isolation.

---

## Project Structure

```plain
kfn-app/
├── cmd/                   # Entry points for each component
│   ├── faasctl/           # CLI tool
│   ├── scheduler/         # Main scheduler/orchestrator
│   ├── api-gateway/       # HTTP API gateway
│   └── runtime/           # Runtime wrapper helpers
├── pkg/                   # Shared Go packages
│   ├── registry/          # etcd client, function metadata
│   ├── runtime/           # Container runtime (containerd/gVisor)
│   ├── build/             # Image build helpers (buildah)
│   ├── events/            # Event bus clients (Kafka/NATS)
│   ├── auth/              # JWT, RBAC, mTLS
│   ├── observability/     # Prometheus metrics, logging
│   └── secrets/           # Vault integration
├── runtime-wrappers/      # Runtime wrapper templates
│   ├── node/              # Node.js wrapper scripts
│   └── python/            # Python wrapper scripts
├── configs/               # Configuration files
├── deploy/                # Kubernetes/helm charts
├── examples/              # Sample functions (Node, Python)
└── tests/                 # Integration tests
```

---

## Phase 1: Foundation & Registry (Weeks 1-4)

### Goal

Set up project structure, Go modules, and the etcd-based function registry.

### Tasks

#### 1.1 Project Bootstrap

- [x] Initialize Go module (`go mod init github.com/geekzy/kfn-app`)
- [x] Set up directory structure as above
- [x] Create Makefile with build targets
- [x] Set up CI/CD pipeline stub (GitHub Actions or GitLab CI)

#### 1.2 Configuration System

- [x] Create config package with Viper integration
- [x] Define config schema:
  - `registry/etcd-endpoints`
  - `container/containerd-socket`
  - `sandbox/type` (gVisor or Kata)
  - `image-registry/harbor-url`
  - `observability/prometheus-url`
- [x] Create default config file in `configs/default.yaml`

#### 1.3 Etcd Registry (pkg/registry)

- [x] Implement etcd client wrapper
  - `FunctionMetadata` struct with fields: `name`, `language`, `runtimeVersion`, `runtimeImage`, `memory`, `timeout`, `handler`
  - CRUD operations: `CreateFunction()`, `GetFunction()`, `ListFunctions()`, `UpdateFunction()`, `DeleteFunction()`
- [x] Implement function state management: `pending`, `ready`, `error`, `deleting`
- [x] Add unit tests with mock etcd

#### 1.4 CLI Skeleton (cmd/faasctl)

- [x] Set up Cobra-based CLI with subcommands:
  - `faasctl deploy`
  - `faasctl list`
  - `faasctl invoke`
  - `faasctl delete`
  - `faasctl status`
- [x] Implement `list` command (reads from etcd via registry package)
- [x] Implement `status` command (shows function state)

---

## Phase 2: Scheduler & Container Runtime (Weeks 5-8)

### Goal

Implement the core scheduler that launches and manages function containers.

### Tasks

#### 2.1 Container Runtime Package (pkg/runtime)

- [ ] Create containerd client wrapper
- [ ] Implement `ContainerSpec` struct defining container configuration
- [ ] Implement container lifecycle:
  - `CreateContainer()` - creates container with gVisor runtime
  - `StartContainer()` - starts the container
  - `StopContainer()` - stops container
  - `RemoveContainer()` - removes container
- [ ] Implement warm pool management:
  - Pre-create containers for popular functions
  - Container reuse logic
  - Idle timeout handling

#### 2.2 Scheduler Core (cmd/scheduler)

- [ ] Create main scheduler service
- [ ] Implement scheduling logic:
  - Function invocation request queue (in-memory or Kafka-backed)
  - Container selection (warm pool or create new)
  - Resource awareness (CPU/memory limits)
- [ ] Implement autoscaling:
  - Scale up based on request queue depth
  - Scale down based on idle container timeout
- [ ] Add Prometheus metrics: `faas_scheduler_invocations`, `faas_container_count`, `faas_cold_starts`

#### 2.3 gVisor Integration

- [ ] Configure containerd to use gVisor as OCI runtime
- [ ] Test container creation with gVisor sandbox
- [ ] Verify isolation (disable SYS_PTRACE, CAP_SYS_ADMIN)

#### 2.4 API Gateway (cmd/api-gateway)

- [ ] Create HTTP API gateway
- [ ] Implement endpoints:
  - `POST /v1/functions/{name}/invoke` - invokes function
  - `GET /v1/functions` - lists functions
  - `GET /v1/functions/{name}` - gets function metadata
  - `POST /v1/functions` - creates function (alternative to CLI)
- [ ] Implement request routing to scheduler
- [ ] Add request ID generation and tracing headers

---

## Phase 3: Runtime Wrappers & Image Building (Weeks 9-12)

### Goal

Create runtime wrappers for Node.js and Python, and implement image build system.

### Tasks

#### 3.1 Node.js Runtime Wrapper (runtime-wrappers/node/)

- [ ] Create entrypoint script `entrypoint.sh`:
  - Install dependencies if `node_modules` missing
  - Start user code with `node index.js`
  - Handle signals gracefully
- [ ] Create HTTP server wrapper `wrapper.js`:
  - Expose `/health` endpoint
  - Expose `/invoke` endpoint for function calls
  - Transform requests to Lambda-style event format
  - Handle timeouts
- [ ] Create base Dockerfile for Node images

#### 3.2 Python Runtime Wrapper (runtime-wrappers/python/)

- [ ] Create entrypoint script `entrypoint.sh`:
  - Install dependencies from `requirements.txt`
  - Start user code with `uvicorn`
- [ ] Create FastAPI wrapper `wrapper.py`:
  - Expose `/health` endpoint
  - Expose `/invoke` endpoint
  - Transform requests to event format
- [ ] Create base Dockerfile for Python images

#### 3.3 Image Build System (pkg/build)

- [ ] Implement buildah integration:
  - `BuildImage()` function that builds OCI images
  - Dynamic Dockerfile generation based on language
  - Image tagging with timestamp and version
- [ ] Implement Harbor registry push:
  - `PushImage()` function
  - Authentication handling
- [ ] Implement image signing with cosign:
  - `SignImage()` function
  - `VerifyImage()` function

#### 3.4 CLI - Deploy Command (cmd/faasctl)

- [ ] Implement `faasctl deploy`:
  - Read source from directory or zip file
  - Validate language and runtime version
  - Call image build system
  - Sign and push to Harbor
  - Register function in etcd

---

## Phase 4: Eventing & Triggers (Weeks 13-16)

### Goal

Add event-driven triggers using Kafka and NATS.

### Tasks

#### 4.1 Event Bus Package (pkg/events)

- [ ] Implement Kafka consumer:
  - Consume from configurable topics
  - Decode events
  - Forward to scheduler
- [ ] Implement NATS JetStream consumer (alternative)
- [ ] Implement event publisher for audit logging

#### 4.2 Trigger System (pkg/triggers)

- [ ] Define `Trigger` struct: `type`, `source`, `function`, `config`
- [ ] Implement trigger storage in etcd
- [ ] Implement trigger management CLI commands:
  - `faasctl trigger create`
  - `faasctl trigger list`
  - `faasctl trigger delete`
- [ ] Implement trigger activation/deactivation

#### 4.3 Kafka Consumer Service

- [ ] Create dedicated Kafka consumer service
- [ ] Implement topic subscription based on triggers
- [ ] Map events to function invocations

#### 4.4 Audit Logging

- [ ] Log all deploy/remove/invocation actions to Kafka audit topic
- [ ] Create audit log viewer CLI command

---

## Phase 5: Security & Auth (Weeks 17-20)

### Goal

Implement authentication, authorization, and image signing verification.

### Tasks

#### 5.1 JWT Authentication (pkg/auth)

- [ ] Implement JWT generation for CLI/API clients
- [ ] Implement JWT validation middleware for API gateway
- [ ] Define scopes: `faas:deploy`, `faas:invoke`, `faas:delete`, `faas:list`
- [ ] Implement token refresh logic

#### 5.2 RBAC System

- [ ] Define tenant/user roles
- [ ] Implement permission checking
- [ ] Add tenant isolation in function listing/invocation

#### 5.3 mTLS

- [ ] Generate CA and certificates for internal services
- [ ] Implement mTLS server/client for scheduler and API gateway
- [ ] Certificate rotation logic

#### 5.4 Image Verification

- [ ] Implement cosign verification in scheduler before container creation
- [ ] Add policy to reject unsigned images
- [ ] Add configurable trust store

#### 5.5 Secret Management (pkg/secrets)

- [ ] Implement Vault client
- [ ] Implement secret storage per function
- [ ] Implement secret injection at container runtime (environment variables)

---

## Phase 6: Observability (Weeks 21-24)

### Goal

Add comprehensive metrics, logging, and dashboards.

### Tasks

#### 6.1 Metrics Package (pkg/observability)

- [ ] Implement Prometheus metrics collector:
  - `faas_invocations_total` (by function)
  - `faas_invocation_duration_seconds` (histogram)
  - `faas_cold_starts_total`
  - `faas_errors_total` (by type)
  - `faas_active_containers`
- [ ] Add metrics to all critical paths

#### 6.2 Logging

- [ ] Implement structured JSON logging using `sirupsen/logrus` or `uber-go/zap`
- [ ] Add request ID correlation
- [ ] Implement log forwarding to Loki

#### 6.3 Grafana Dashboards

- [ ] Create dashboard for function overview
- [ ] Create dashboard for system health
- [ ] Create dashboard for performance metrics

#### 6.4 Distributed Tracing (Optional)

- [ ] Add OpenTelemetry integration
- [ ] Trace requests from API gateway through scheduler to function

---

## Phase 7: Testing & CI/CD (Weeks 25-28)

### Goal
Comprehensive test coverage and automated CI/CD pipeline.

### Tasks

#### 7.1 Unit Tests

- [ ] Achieve >80% code coverage
- [ ] Add unit tests for all packages
- [ ] Use `golang/mock` for mocking dependencies

#### 7.2 Integration Tests

- [ ] Create test environment with all services
- [ ] Test full deploy/invoke cycle
- [ ] Test error scenarios
- [ ] Test autoscaling behavior

#### 7.3 CI/CD Pipeline

- [ ] Set up automated builds
- [ ] Set up automated tests
- [ ] Set up image promotion (dev -> staging -> prod)
- [ ] Set up security scanning

#### 7.4 Load Testing

- [ ] Create load test scripts (k6 or locust)
- [ ] Target: <50ms latency at 99th percentile
- [ ] Target: >1000 concurrent requests

---

## Phase 8: Documentation & Demo (Weeks 29-32)

### Goal

Complete documentation and demo scenarios.

### Tasks

#### 8.1 Documentation

- [ ] README with quick start
- [ ] CLI reference documentation
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Architecture documentation
- [ ] Deployment guide
- [ ] Troubleshooting guide

#### 8.2 Example Functions

- [ ] Create Node.js hello world example
- [ ] Create Python hello world example
- [ ] Create event-driven examples
- [ ] Create multi-language demo

#### 8.3 Demo Scripts

- [ ] Create live demo script
- [ ] Create video walkthrough
- [ ] Create workshop materials

---

## Critical Files by Phase

### Phase 1 (Foundation)

```plain
go.mod
Makefile
configs/default.yaml
pkg/registry/client.go
pkg/registry/metadata.go
cmd/faasctl/main.go
cmd/faasctl/list.go
cmd/faasctl/status.go
```

### Phase 2 (Scheduler)

```plain
pkg/runtime/containerd.go
pkg/runtime/container.go
pkg/scheduler/scheduler.go
cmd/scheduler/main.go
cmd/api-gateway/main.go
```

### Phase 3 (Runtime & Build)

```plain
runtime-wrappers/node/entrypoint.sh
runtime-wrappers/node/wrapper.js
runtime-wrappers/python/entrypoint.sh
runtime-wrappers/python/wrapper.py
pkg/build/buildah.go
pkg/build/registry.go
cmd/faasctl/deploy.go
```

### Phase 4 (Eventing)

```plain
pkg/events/kafka.go
pkg/events/nats.go
pkg/triggers/trigger.go
cmd/scheduler/consumer.go
```

### Phase 5 (Security)

```plain
pkg/auth/jwt.go
pkg/auth/rbac.go
pkg/secrets/vault.go
```

### Phase 6 (Observability)

```plain
pkg/observability/metrics.go
pkg/observability/logging.go
configs/grafana/dashboards/
```

---

## Success Criteria

| Criterion | Target |
| ----------- | -------- |
| Node.js function deployment | Working |
| Python function deployment | Working |
| Function invocation via CLI | Working |
| Function invocation via HTTP API | Working |
| Cold start latency | <100ms |
| Warm invocation latency | <50ms at 99th percentile |
| Security (image signing, sandboxing) | Passed audit |
| Autoscaling | Working |
| Observability (metrics, logs) | Working |
| Documentation | Complete |
| End-to-end demo | Successful |
