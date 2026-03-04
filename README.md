# KFN FaaS Platform

The KFN (Kubernetes Function-as-a-Service) platform is a lightweight, Go‑based FaaS system that supports Node.js and Python functions with container‑per‑invocation isolation. It aims to provide a simple, highly‑scalable runtime for serverless workloads running on Kubernetes or locally.

## Key Features

- **Multi‑language support** – Run Node.js and Python functions out of the box using lightweight Docker images.
- **Container‑per‑invocation** – Each request is isolated in its own container (gVisor or Kata) for strong security.
- **Zero‑trust deployment** – Functions are built, signed, and pushed to a trusted registry automatically.
- **Observability** – Built‑in Prometheus metrics, structured logging, and event tracing.
- **CLI & API** – `faasctl` CLI for everyday operations and a REST API gateway for programmatic access.
- **Custom runtime wrappers** – Extendable templates for additional languages.

## Getting Started

```bash
# Install dependencies
make deps

# Build the CLI
make build

# Deploy a function (e.g., a Node.js hello world)
faasctl deploy ./examples/hello-node
```

For a full walkthrough, see the [Documentation](docs/).