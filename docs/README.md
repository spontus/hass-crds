# hass-crds Documentation

**hass-crds** is a Kubernetes controller that manages Home Assistant entities via MQTT autodiscovery. Define your HA entities as Kubernetes Custom Resources, and the controller automatically publishes MQTT discovery payloads so Home Assistant creates and manages them for you -- fully declarative, GitOps-friendly.

## Prerequisites

- Kubernetes 1.24+
- An MQTT broker accessible from the cluster (e.g. Mosquitto)
- Home Assistant with [MQTT integration](https://www.home-assistant.io/integrations/mqtt/) and discovery enabled
- `kubectl` configured for your cluster

## Documentation

| Section | Description |
|---|---|
| [Architecture](architecture.md) | Design decisions, API group, reconciliation model |
| [Getting Started](getting-started.md) | Install CRDs, deploy the controller, verify |
| [Controller](controller.md) | MQTT connection, reconciliation loop, deletion, status |
| [CRD Reference](crds/README.md) | All supported entity types and their fields |
| [Common Fields](crds/common-fields.md) | Shared fields across all CRD types |
| [Examples](examples/README.md) | Copy-pasteable YAML manifests for every entity type |
| [Admission Webhooks](admission-webhooks.md) | Validating webhook for cross-resource checks |
| [Entity Type Coverage](future-entity-types.md) | All 28 supported entity types |

## Supported Entity Types

All 28 Home Assistant MQTT entity types are supported, plus a utility `MQTTDevice` CRD for shared device definitions. See [CRD Reference](crds/README.md) for the full list.

## Quick Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: restart-server
spec:
  name: "Restart Server"
  commandTopic: "homeassistant/button/restart-server/command"
  device:
    name: "My Kubernetes Cluster"
    identifiers:
      - "k8s-cluster-01"
```

Apply it:

```bash
kubectl apply -f button.yaml
```

The controller publishes a discovery payload to MQTT, and Home Assistant automatically creates a button entity.
