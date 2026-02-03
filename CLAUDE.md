# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

hass-crds is a Kubernetes operator that manages Home Assistant entities via MQTT autodiscovery. Users define HA entities as Kubernetes Custom Resources, and the controller publishes MQTT discovery payloads so Home Assistant automatically creates them.

- **API Group**: `mqtt.home-assistant.io/v1alpha1`
- **Language**: Go 1.22 with controller-runtime
- **CRD Generation**: Python scripts (not kubebuilder markers)

## Common Commands

```bash
# Build
make build                    # Build controller binary
make generate                 # Generate DeepCopy methods
make manifests                # Generate RBAC manifests

# CRD Generation (Python-based, not controller-gen)
make crds                     # Generate CRDs via Python
cd config/crd/generator && python generate.py

# Test
make test                     # Unit tests with envtest
go test ./internal/payload/...  # Run single package tests
make test-e2e                 # E2E tests (requires Kind cluster)

# Lint
make lint                     # Run golangci-lint
make lint-fix                 # Auto-fix lint issues

# Deploy
make install                  # Install CRDs to cluster
make deploy                   # Deploy controller
make run                      # Run controller locally
```

## Architecture

### Key Components

1. **API Types** (`api/v1alpha1/`)
   - One CRD per HA entity type (MQTTButton, MQTTSensor, MQTTLight, etc.)
   - `common_types.go` defines shared structs: `CommonSpec`, `CommonStatus`, `DeviceBlock`, `AvailabilityConfig`
   - All entity types embed `CommonSpec` and implement `GetCommonSpec()`/`GetCommonStatus()` interfaces

2. **Controllers** (`internal/controller/`)
   - `base.go`: `BaseReconciler` with shared reconciliation logic (publish discovery, handle deletion, manage finalizers)
   - Entity-specific controllers call `BaseReconciler.PublishDiscovery()` with a payload builder function
   - Namespace-scoped (each controller watches only its own namespace)

3. **Payload Builder** (`internal/payload/`)
   - Converts CRD specs to MQTT discovery JSON
   - Auto-converts camelCase to snake_case for HA compatibility
   - `DeviceBlockToMap()` and `AvailabilityToMap()` for nested structures

4. **Topic Generation** (`internal/topic/`)
   - Discovery topic: `<prefix>/<component>/<namespace>-<name>/config`
   - Auto-generates `unique_id` from `<namespace>-<name>` if not specified

5. **CRD Generator** (`config/crd/generator/`)
   - Python-based CRD generation (not kubebuilder markers)
   - Entity schemas in `schemas/entities.py`
   - Common fields in `schemas/common.py`
   - Run `make crds` after modifying schemas

### Reconciliation Flow

```
CRD Create/Update → Add finalizer → Build payload → Publish to MQTT → Update status
CRD Delete → Publish empty payload → Remove finalizer → Allow deletion
```

### Field Naming Convention

CRD specs use camelCase (K8s convention), automatically converted to snake_case (HA convention):
- `commandTopic` → `command_topic`
- `deviceClass` → `device_class`

## Adding a New Entity Type

1. Add type definitions in `api/v1alpha1/mqtt<entity>_types.go`
2. Add entity schema in `config/crd/generator/schemas/entities.py`
3. Run `make crds` and `make generate`
4. Add controller in `internal/controller/mqtt<entity>_controller.go`
5. Register controller in `internal/controller/setup.go`

## Testing

- Unit tests: `go test ./... -short`
- Integration tests: `go test ./... -tags=integration` (requires MQTT broker)
- E2E tests: `make test-e2e` (requires Docker, Kind, kubectl)
- Payload builder tests: `internal/payload/builder_test.go`

### E2E Tests with Kind

The e2e tests use a Kind cluster with Mosquitto MQTT broker, Home Assistant, and the controller deployed. Tests verify MQTT discovery payloads are published correctly.

**Prerequisites**: Docker, Kind, kubectl

#### Running E2E Tests

```bash
# Full cycle: create cluster → run tests → delete cluster
make test-e2e

# Run tests but keep cluster running (for debugging failures)
make test-e2e-keep

# Reuse existing cluster (fast iteration after test-e2e-keep)
make test-e2e-reuse
```

#### Manual Cluster Management

```bash
# Create cluster without running tests
make e2e-cluster-create

# Delete cluster
make e2e-cluster-delete

# Access cluster after creation
kubectl --context kind-hass-crds-e2e get pods -n hass-crds-e2e
```

#### Environment Variables

| Variable | Effect |
|----------|--------|
| `SKIP_CLUSTER_SETUP=true` | Skip Kind cluster creation (reuse existing) |
| `SKIP_CLUSTER_TEARDOWN=true` | Keep cluster after tests complete |

#### Port Mappings (for debugging)

| Service | Host Port | Description |
|---------|-----------|-------------|
| Home Assistant | `localhost:8123` | HA web UI |
| Mosquitto MQTT | `localhost:1883` | MQTT broker |

#### Test Structure

- `test/e2e/e2e_suite_test.go` - Ginkgo suite setup (cluster lifecycle)
- `test/e2e/e2e_test.go` - Test cases (MQTTButton, MQTTSensor, MQTTSwitch)
- `test/e2e/kind-config.yaml` - Kind cluster configuration
- `test/e2e/manifests/` - K8s manifests for Mosquitto, HA, controller
- `test/utils/utils.go` - Helper functions (MQTT subscribe, pod logs, etc.)

#### Debugging Tips

```bash
# View controller logs
kubectl logs -n hass-crds-e2e deployment/hass-crds-controller -f

# View Mosquitto logs
kubectl logs -n hass-crds-e2e deployment/mosquitto -f

# Subscribe to MQTT topic from inside cluster
kubectl exec -n hass-crds-e2e deployment/mosquitto -- \
  mosquitto_sub -h localhost -t "homeassistant/#" -v

# Apply test resource manually
kubectl apply -n hass-crds-e2e -f - <<EOF
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: debug-button
spec:
  name: "Debug Button"
  commandTopic: "test/button/command"
EOF
```

## Linting

Uses golangci-lint with these enabled linters: errcheck, govet, staticcheck, gofmt, goimports, misspell, unused, ineffassign, goconst, gocyclo.
