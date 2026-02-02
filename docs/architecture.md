# Architecture

## Problem Statement

Home Assistant supports [MQTT discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery) -- a mechanism where a device publishes a JSON configuration payload to a well-known MQTT topic, and HA automatically creates the corresponding entity. This is powerful, but managing discovery payloads manually is tedious, error-prone, and doesn't fit well into infrastructure-as-code workflows.

**hass-crds** solves this by letting you define HA entities as Kubernetes Custom Resources. A controller running in your cluster watches these resources and publishes the appropriate MQTT discovery payloads. This makes entity management:

- **Declarative** -- define what you want, not how to get there
- **GitOps-friendly** -- store entity definitions in version control
- **Validated** -- CRD schemas catch errors before they reach MQTT
- **Observable** -- status subresources show what's been published and when

## High-Level Flow

```
User creates CRD instance (kubectl apply / GitOps)
        |
        v
Controller watches CRD instances in its namespace
        |
        v
Controller builds MQTT discovery JSON from CRD spec
  - camelCase fields -> snake_case keys
  - auto-generates unique_id if not set
  - derives discovery topic from component/namespace/name
        |
        v
Controller publishes JSON to MQTT with retain=true
        |
        v
Home Assistant receives discovery payload, creates entity
        |
        v
Controller updates CRD status (.status.lastPublished, .status.discoveryTopic)
```

## API Group and Versioning

- **API group**: `mqtt.home-assistant.io`
- **Version**: `v1alpha1`
- **Full apiVersion**: `mqtt.home-assistant.io/v1alpha1`

The `v1alpha1` version indicates this API is under active development. Breaking changes may occur before `v1beta1`.

## CRD Design Decisions

### One CRD per Entity Type

Each Home Assistant entity type (button, switch, light, etc.) gets its own CRD rather than a generic `MQTTEntity` kind. This provides:

- **Schema validation** -- each CRD validates only the fields relevant to that entity type
- **Clear documentation** -- each type has its own reference page
- **Better kubectl experience** -- `kubectl get mqttbuttons` is more useful than `kubectl get mqttentities --field-selector type=button`

### CRD Naming Convention

| Component | Kind | Resource (plural) | Singular |
|---|---|---|---|
| `button` | `MQTTButton` | `mqttbuttons` | `mqttbutton` |
| `switch` | `MQTTSwitch` | `mqttswitches` | `mqttswitch` |
| `light` | `MQTTLight` | `mqttlights` | `mqttlight` |
| `sensor` | `MQTTSensor` | `mqttsensors` | `mqttsensor` |
| `binary_sensor` | `MQTTBinarySensor` | `mqttbinarysensors` | `mqttbinarysensor` |
| `cover` | `MQTTCover` | `mqttcovers` | `mqttcover` |
| `climate` | `MQTTClimate` | `mqttclimates` | `mqttclimate` |
| `number` | `MQTTNumber` | `mqttnumbers` | `mqttnumber` |
| `select` | `MQTTSelect` | `mqttselects` | `mqttselect` |
| `text` | `MQTTText` | `mqtttexts` | `mqtttext` |

### Field Naming: camelCase to snake_case

CRD specs use **camelCase** (Kubernetes convention), which the controller maps to **snake_case** (Home Assistant MQTT convention) when building the discovery JSON:

```
CRD spec field        MQTT JSON key
──────────────        ─────────────
commandTopic     ->   command_topic
stateTopic       ->   state_topic
uniqueId         ->   unique_id
deviceClass      ->   device_class
```

### Auto-generated `unique_id`

If `spec.uniqueId` is not set, the controller generates one from `<namespace>-<name>`. This ensures every entity has a stable unique ID that HA can use for entity registry tracking.

### Discovery Topic Derivation

The controller publishes discovery payloads to:

```
<discovery_prefix>/<component>/<namespace>-<name>/config
```

Where:

- `discovery_prefix` defaults to `homeassistant` (configurable via `MQTT_DISCOVERY_PREFIX`)
- `component` is the HA component type (e.g. `button`, `switch`, `light`)
- `namespace` and `name` come from the CRD instance's Kubernetes metadata

Example: An `MQTTButton` named `restart-server` in the `default` namespace publishes to:

```
homeassistant/button/default-restart-server/config
```

### Retain Behavior

- **Discovery messages** are always published with `retain=true` so HA picks them up even after a restart
- The `retain` field in a CRD's spec controls whether the entity's command/state topics use retained messages -- this is separate from discovery retention

### Finalizers for Clean Deletion

When a CRD instance is deleted, the controller:

1. Intercepts deletion via a finalizer (`mqtt.home-assistant.io/discovery`)
2. Publishes an **empty payload** to the discovery topic (this tells HA to remove the entity)
3. Removes the finalizer, allowing Kubernetes to complete the deletion

### Namespace-Scoped

The controller watches only its own namespace. This provides:

- **Isolation** -- multiple controllers in different namespaces don't interfere
- **Simpler RBAC** -- no cluster-wide permissions needed
- **Multi-tenancy** -- different teams can manage their own HA entities

### MQTTDevice: Shared Device Definitions

Multiple entities often share the same device block (e.g. a sensor hub with temperature, humidity, and pressure sensors). Rather than duplicating the device block on each entity, a dedicated `MQTTDevice` CRD holds the device metadata, and entities reference it via `deviceRef`:

```yaml
spec:
  deviceRef:
    name: "weather-station"
```

This keeps entity CRDs focused on their type-specific fields and ensures device metadata is consistent across entities. See [MQTTDevice](crds/device.md) for details.

### Secret References

Sensitive values (alarm codes, lock codes) should not be stored in plaintext in CRD specs. Any string field can use a `secretRef` to load its value from a Kubernetes Secret at reconciliation time:

```yaml
spec:
  codeFormat:
    secretRef:
      name: "lock-codes"
      key: "front-door-code"
```

See [Common Fields — Secret References](crds/common-fields.md#secret-references) for details.

### Admission Webhooks

An optional validating webhook catches errors at apply time that CRD schema validation cannot:

- Duplicate `uniqueId` values across CRs in a namespace
- Duplicate discovery topics
- Invalid field combinations (e.g. wrong fields for a light schema)
- Missing referenced `MQTTDevice` or Secret resources

See [Admission Webhooks](admission-webhooks.md) for setup and the full list of validation rules.

### Common Fields as Convention

Kubernetes CRDs don't support inheritance, so shared fields (device, availability, entity metadata, MQTT options) are **embedded in each CRD type** and documented centrally in [common-fields.md](crds/common-fields.md). Each entity type's CRD includes the full set of common fields alongside its type-specific fields.

### Multi-Tenancy via Namespaces

The controller is namespace-scoped by design. Each namespace with its own controller deployment acts as an independent tenant with its own MQTT broker connection, discovery prefix, and topic prefix. Multiple tenants can share a single broker by using distinct prefixes. See [Controller — Multi-Tenancy](controller.md#multi-tenancy) for deployment patterns.
