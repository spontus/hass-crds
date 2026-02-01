# Common Fields

All CRD types embed the following shared field groups. These are documented here once; each entity type's reference page links back to this document.

## Entity Metadata

Basic fields present on every entity type.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `name` | `name` | `string` | No | Resource name | Display name in HA |
| `uniqueId` | `unique_id` | `string` | No | `<namespace>-<name>` | Unique identifier for HA entity registry |
| `icon` | `icon` | `string` | No | -- | MDI icon (e.g. `mdi:thermometer`) |
| `entityCategory` | `entity_category` | `string` | No | -- | `config` or `diagnostic` |
| `enabledByDefault` | `enabled_by_default` | `bool` | No | `true` | Whether the entity is enabled when first discovered |
| `objectId` | `object_id` | `string` | No | -- | Override for HA entity ID generation |

## Device

Associates the entity with a device in HA's device registry. Multiple entities can share the same device block to be grouped together.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `device.name` | `device.name` | `string` | No | -- | Device display name |
| `device.identifiers` | `device.identifiers` | `[]string` | No | -- | List of identifiers (at least one of `identifiers` or `connections` is needed) |
| `device.connections` | `device.connections` | `[][]string` | No | -- | List of `[type, value]` pairs (e.g. `[["mac", "aa:bb:cc:dd:ee:ff"]]`) |
| `device.manufacturer` | `device.manufacturer` | `string` | No | -- | Device manufacturer |
| `device.model` | `device.model` | `string` | No | -- | Device model |
| `device.modelId` | `device.model_id` | `string` | No | -- | Device model identifier |
| `device.serialNumber` | `device.serial_number` | `string` | No | -- | Device serial number |
| `device.hwVersion` | `device.hw_version` | `string` | No | -- | Hardware version |
| `device.swVersion` | `device.sw_version` | `string` | No | -- | Software version |
| `device.suggestedArea` | `device.suggested_area` | `string` | No | -- | Suggested area in HA (e.g. "Living Room") |
| `device.configurationUrl` | `device.configuration_url` | `string` | No | -- | URL for device configuration |
| `device.viaDevice` | `device.via_device` | `string` | No | -- | Identifier of device that routes messages |

### Example

```yaml
spec:
  device:
    name: "Living Room Sensor Hub"
    identifiers:
      - "sensor-hub-01"
    manufacturer: "Custom"
    model: "SensorHub v2"
    swVersion: "1.4.0"
    suggestedArea: "Living Room"
```

### Device Reference

Instead of duplicating the device block on every entity, you can create a shared [`MQTTDevice`](device.md) resource and reference it by name:

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `deviceRef.name` | -- | `string` | No | -- | Name of an `MQTTDevice` resource in the same namespace |

When `deviceRef` is set, the controller resolves the referenced `MQTTDevice` and injects its fields into the discovery payload. If both `device` and `deviceRef` are present, the inline `device` block takes priority.

```yaml
spec:
  deviceRef:
    name: "weather-station"
```

See [`MQTTDevice`](device.md) for details and a full example.

## Availability

Defines how HA determines whether the entity is available. If omitted, HA uses the default MQTT availability topic.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `availability` | `availability` | `[]object` | No | -- | List of availability topics |
| `availability[].topic` | `availability[].topic` | `string` | Yes | -- | MQTT topic for availability |
| `availability[].payloadAvailable` | `availability[].payload_available` | `string` | No | `online` | Payload indicating available |
| `availability[].payloadNotAvailable` | `availability[].payload_not_available` | `string` | No | `offline` | Payload indicating unavailable |
| `availability[].valueTemplate` | `availability[].value_template` | `string` | No | -- | Template to extract availability from payload |
| `availabilityMode` | `availability_mode` | `string` | No | `latest` | `all`, `any`, or `latest` |

### Example

```yaml
spec:
  availability:
    - topic: "devices/sensor-hub-01/status"
      payloadAvailable: "online"
      payloadNotAvailable: "offline"
  availabilityMode: "all"
```

## MQTT Options

Control MQTT behavior for the entity's command and state topics.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `qos` | `qos` | `integer` | No | `0` | MQTT QoS level (0, 1, or 2) |
| `retain` | `retain` | `bool` | No | `false` | Whether to retain messages on command/state topics |
| `encoding` | `encoding` | `string` | No | `utf-8` | Payload encoding |

> **Note**: The `retain` field here controls the entity's command/state topics. Discovery messages are always published with `retain=true` regardless of this setting.

## JSON Attributes

Allows the entity to expose additional attributes from a JSON payload.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `jsonAttributesTopic` | `json_attributes_topic` | `string` | No | -- | MQTT topic for JSON attributes |
| `jsonAttributesTemplate` | `json_attributes_template` | `string` | No | -- | Template to extract attributes from payload |

## Rediscovery

Controls periodic re-publishing of the MQTT discovery payload. This is useful to ensure Home Assistant picks up entity configurations after restarts, broker failovers, or retained-message expiration.

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `rediscoverInterval` | -- | `duration` | No | -- | How often to re-publish the discovery config payload (e.g. `5m`, `1h`). If omitted, discovery is only published on resource creation or update. |

> **Note**: This is a controller-only field and is not included in the MQTT discovery JSON sent to Home Assistant. The value is a Kubernetes-style duration string.

### Example

```yaml
spec:
  rediscoverInterval: "30m"
```

## Secret References

Sensitive field values (e.g. alarm codes, lock codes) can be loaded from Kubernetes Secrets instead of stored in plaintext in the CRD spec. Any string field in the CRD spec can use a `secretRef` instead of a literal value.

### Syntax

Replace the literal value with a `secretRef` block:

```yaml
spec:
  # Instead of this:
  # codeFormat: "1234"

  # Use this:
  codeFormat:
    secretRef:
      name: "lock-codes"
      key: "front-door-code"
```

| Field | Type | Required | Description |
|---|---|---|---|
| `secretRef.name` | `string` | Yes | Name of the Kubernetes Secret in the same namespace |
| `secretRef.key` | `string` | Yes | Key within the Secret's `data` to read the value from |

### Behavior

- The controller reads the Secret value at reconciliation time and injects it into the MQTT discovery payload
- If the Secret or key does not exist, the controller sets the `Published` condition to `False` with reason `SecretNotFound`
- Changes to the referenced Secret trigger re-reconciliation of all CRDs that reference it
- Secret values are never logged by the controller

### Example

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: alarm-codes
  namespace: hass-crds
type: Opaque
stringData:
  panel-code: "1234"
---
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTAlarmControlPanel
metadata:
  name: home-alarm
  namespace: hass-crds
spec:
  name: "Home Alarm"
  commandTopic: "alarm/home/set"
  stateTopic: "alarm/home/state"
  commandTemplate:
    secretRef:
      name: "alarm-codes"
      key: "panel-code"
  device:
    name: "Alarm Panel"
    identifiers:
      - "alarm-panel-01"
```
