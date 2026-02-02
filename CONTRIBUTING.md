# Contributing to hass-crds

Thank you for your interest in contributing to hass-crds.

## Development Setup

### Prerequisites

- Go 1.21+
- Docker
- kubectl configured for a test cluster
- An MQTT broker (Mosquitto recommended for local dev)
- Home Assistant instance (optional, for end-to-end testing)

### Clone and Build

```bash
git clone git@github.com:spontus/hass-crds.git
cd hass-crds

# Install dependencies
go mod download

# Build the controller
go build -o bin/controller ./cmd/controller

# Run tests
go test ./...
```

### Local Development

1. Start a local MQTT broker:
   ```bash
   docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto
   ```

2. Install CRDs to your cluster:
   ```bash
   kubectl apply -f config/crd/crds.yaml
   ```

3. Run the controller locally:
   ```bash
   export MQTT_HOST=localhost
   export MQTT_PORT=1883
   export LOG_LEVEL=debug
   go run ./cmd/controller
   ```

4. Create a test resource:
   ```bash
   kubectl apply -f docs/examples/basic-button.yaml
   ```

## Adding a New Entity Type

To add support for a new Home Assistant entity type:

### 1. Create the CRD Definition

Add a new file in `config/crd/bases/`:

```yaml
# config/crd/bases/mqtt.home-assistant.io_mqttnewtypes.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: mqttnewtypes.mqtt.home-assistant.io
spec:
  group: mqtt.home-assistant.io
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          # ... define schema
  scope: Namespaced
  names:
    plural: mqttnewtypes
    singular: mqttnewtype
    kind: MQTTNewType
```

### 2. Add Controller Logic

Add reconciliation logic in `internal/controller/`:

```go
// internal/controller/newtype_controller.go
package controller

func (r *NewTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Fetch the resource
    // Build the MQTT discovery payload
    // Publish to MQTT
    // Update status
}
```

### 3. Add Documentation

Create documentation following the existing pattern:

```bash
# CRD reference
docs/crds/new-type.md

# Example manifest
docs/examples/new-type.yaml
```

Use existing files as templates. Include:
- Type-specific fields table with MQTT key mappings
- Link to Home Assistant documentation
- Working example YAML
- Navigation footer linking to related types

### 4. Update Indexes

Update these files to include the new type:
- `docs/crds/README.md` - Add to entity types table
- `docs/future-entity-types.md` - Move from future to implemented
- `docs/getting-started.md` - Add to RBAC rules

### 5. Add Tests

```go
// internal/controller/newtype_controller_test.go
func TestNewTypeReconciler(t *testing.T) {
    // Test creation publishes correct payload
    // Test update re-publishes
    // Test deletion publishes empty payload
    // Test secretRef resolution
}
```

## Code Style

### Go

- Follow standard Go conventions
- Run `go fmt` and `go vet` before committing
- Use meaningful variable names
- Add comments for non-obvious logic

### Documentation

- Use [CommonMark](https://commonmark.org/) markdown
- Keep line length under 120 characters
- Use tables for field references
- Include working examples

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(crd): add MQTTNewType support

- Add CRD definition
- Add controller reconciliation
- Add documentation and examples

Closes #123
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

## Pull Request Process

1. **Fork** the repository
2. **Create a branch**: `git checkout -b feat/my-feature`
3. **Make changes** and commit with conventional commit messages
4. **Test locally**: Run `go test ./...` and verify with a real cluster
5. **Push**: `git push origin feat/my-feature`
6. **Open PR**: Include description of changes and testing done

### PR Checklist

- [ ] Tests pass locally
- [ ] Documentation updated (if applicable)
- [ ] CRD schema validates correctly
- [ ] Example manifests work
- [ ] RBAC rules updated (if new CRD type)

## Testing

### Unit Tests

```bash
go test ./... -v
```

### Integration Tests

Requires a Kubernetes cluster and MQTT broker:

```bash
export KUBECONFIG=~/.kube/config
export MQTT_HOST=localhost
go test ./... -tags=integration -v
```

### End-to-End Tests

Requires Home Assistant:

```bash
export HA_URL=http://localhost:8123
export HA_TOKEN=<long-lived-access-token>
go test ./... -tags=e2e -v
```

## Reporting Issues

### Bug Reports

Include:
- Controller version
- Kubernetes version
- MQTT broker type/version
- Home Assistant version
- CRD manifest (redact secrets)
- Controller logs
- Steps to reproduce

### Feature Requests

Describe:
- Use case
- Proposed solution
- Alternative approaches considered

## Code of Conduct

Be respectful and constructive. We're all here to improve the project.

## License

By contributing, you agree that your contributions will be licensed under the project's license.
