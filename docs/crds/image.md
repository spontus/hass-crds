# MQTTImage

An image entity that displays a static image from an MQTT topic or URL.

- **Kind**: `MQTTImage`
- **Resource**: `mqttimages`
- **HA Component**: `image`
- **HA Docs**: [MQTT Image](https://www.home-assistant.io/integrations/image.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `imageTopic` | `image_topic` | `string` | No | -- | Topic to receive raw image data |
| `imageEncoding` | `image_encoding` | `string` | No | -- | `b64` for base64-encoded images |
| `urlTopic` | `url_topic` | `string` | No | -- | Topic to receive image URL |
| `urlTemplate` | `url_template` | `string` | No | -- | Template to extract URL from payload |
| `contentType` | `content_type` | `string` | No | `image/png` | Image MIME type |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTImage
metadata:
  name: plant-image
spec:
  name: "Plant Cam"
  imageTopic: "cameras/plant/image"
  imageEncoding: "b64"
  contentType: "image/jpeg"
  device:
    name: "Plant Camera"
    identifiers:
      - "plant-cam-01"
    manufacturer: "Custom"
    model: "ESP32-CAM"
    suggestedArea: "Garden"
```
