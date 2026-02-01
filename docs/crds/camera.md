# MQTTCamera

A camera entity that receives images via MQTT.

- **Kind**: `MQTTCamera`
- **Resource**: `mqttcameras`
- **HA Component**: `camera`
- **HA Docs**: [MQTT Camera](https://www.home-assistant.io/integrations/camera.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `topic` | `topic` | `string` | Yes | -- | MQTT topic to subscribe to for image data |
| `imageEncoding` | `image_encoding` | `string` | No | -- | `b64` for base64-encoded images |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTCamera
metadata:
  name: doorbell-camera
spec:
  name: "Doorbell Camera"
  topic: "cameras/doorbell/image"
  imageEncoding: "b64"
  device:
    name: "Doorbell"
    identifiers:
      - "doorbell-01"
    manufacturer: "Custom"
    model: "Smart Doorbell v2"
    suggestedArea: "Front Door"
```
