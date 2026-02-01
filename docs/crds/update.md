# MQTTUpdate

A firmware/software update entity that tracks available updates via MQTT.

- **Kind**: `MQTTUpdate`
- **Resource**: `mqttupdates`
- **HA Component**: `update`
- **HA Docs**: [MQTT Update](https://www.home-assistant.io/integrations/update.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic with JSON payload containing update info |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `commandTopic` | `command_topic` | `string` | No | -- | Topic to trigger update installation |
| `payloadInstall` | `payload_install` | `string` | No | `INSTALL` | Payload to trigger installation |
| `latestVersionTopic` | `latest_version_topic` | `string` | No | -- | Topic to read latest available version |
| `latestVersionTemplate` | `latest_version_template` | `string` | No | -- | Template to extract latest version |
| `deviceClass` | `device_class` | `string` | No | -- | `firmware` |
| `entityPicture` | `entity_picture` | `string` | No | -- | URL to an image for the update entity |
| `releaseUrl` | `release_url` | `string` | No | -- | URL to release notes |
| `releaseSummary` | `release_summary` | `string` | No | -- | Summary of the release |
| `title` | `title` | `string` | No | -- | Title of the software/firmware |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTUpdate
metadata:
  name: firmware-update
spec:
  name: "Firmware Update"
  stateTopic: "device/firmware/state"
  latestVersionTopic: "device/firmware/latest"
  commandTopic: "device/firmware/install"
  deviceClass: "firmware"
  title: "Device Firmware"
  releaseUrl: "https://github.com/example/device/releases"
  device:
    name: "IoT Device"
    identifiers:
      - "iot-device-01"
    manufacturer: "Custom"
    model: "IoT Hub v3"
```
