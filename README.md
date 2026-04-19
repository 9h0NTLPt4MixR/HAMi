# HAMi - Heterogeneous AI Computing Virtualization Middleware

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/HAMi-io/HAMi)](https://goreportcard.com/report/github.com/HAMi-io/HAMi)

HAMi is a Kubernetes device plugin and scheduler extension that enables sharing and virtualization of heterogeneous AI accelerators (GPUs, NPUs, etc.) across pods.

## Features

- **GPU Sharing**: Share a single GPU across multiple pods with resource isolation
- **GPU Memory Virtualization**: Set precise GPU memory limits per container
- **Multi-vendor Support**: NVIDIA GPUs, Cambricon MLUs, Hygon DCUs, and more
- **Kubernetes Native**: Works as a standard device plugin and scheduler extender
- **Resource Monitoring**: Built-in metrics and monitoring support

## Prerequisites

- Kubernetes >= 1.23
- Go >= 1.21
- Docker or containerd runtime
- NVIDIA drivers (for GPU support)

## Quick Start

### Installation via Helm

```bash
helm repo add hami https://hami-io.github.io/HAMi
helm repo update
helm install hami hami/hami --namespace kube-system
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/HAMi-io/HAMi.git
cd HAMi

# Build all components
make build

# Build Docker images
make docker-build
```

## Architecture

```
┌─────────────────────────────────────────┐
│              Kubernetes API             │
└──────────────┬──────────────────────────┘
               │
   ┌───────────▼───────────┐
   │   HAMi Scheduler      │  (Scheduler Extender)
   │   Extender            │
   └───────────┬───────────┘
               │
   ┌───────────▼───────────┐
   │   HAMi Device Plugin  │  (Per Node DaemonSet)
   │   (Node Agent)        │
   └───────────────────────┘
```

## Configuration

See [docs/configuration.md](docs/configuration.md) for detailed configuration options.

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) and check the [good first issues](https://github.com/HAMi-io/HAMi/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

### Development Setup

```bash
# Run tests
make test

# Run linter
make lint

# Generate code
make generate
```

## Local Notes

> **Personal fork** — I'm using this to experiment with GPU memory limits on a single-node k3s cluster.
> Main branch tracks upstream; `dev` branch has my local patches.
>
> **k3s install tip**: k3s doesn't use `/etc/docker/daemon.json` — make sure to configure the nvidia container runtime
> via `/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl` instead, otherwise the device plugin won't work.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project is a fork of [Project-HAMi/HAMi](https://github.com/Project-HAMi/HAMi). We thank all the original contributors.
