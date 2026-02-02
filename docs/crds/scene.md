# MQTTScene

A scene entity that can be activated via MQTT.

- **Kind**: `MQTTScene`
- **Resource**: `mqttscenes`
- **HA Component**: `scene`
- **HA Docs**: [MQTT Scene](https://www.home-assistant.io/integrations/scene.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish when scene is activated |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload sent when scene is activated |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTScene
metadata:
  name: movie-night-scene
spec:
  name: "Movie Night"
  commandTopic: "scenes/movie-night/activate"
  payloadOn: "ON"
  icon: "mdi:movie-open"
  device:
    name: "Scene Controller"
    identifiers:
      - "scene-controller-01"
    manufacturer: "Custom"
    model: "Scene Manager"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTCamera](camera.md), [MQTTImage](image.md)
