# MQTTDeviceTracker

A device tracker entity for presence detection and location tracking via MQTT.

- **Kind**: `MQTTDeviceTracker`
- **Resource**: `mqttdevicetrackers`
- **HA Component**: `device_tracker`
- **HA Docs**: [MQTT Device Tracker](https://www.home-assistant.io/integrations/device_tracker.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic to read tracker state (home/not_home or zone name) |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadHome` | `payload_home` | `string` | No | `home` | Payload representing home |
| `payloadNotHome` | `payload_not_home` | `string` | No | `not_home` | Payload representing not home |
| `payloadReset` | `payload_reset` | `string` | No | -- | Payload that resets the tracker to unknown |
| `sourceType` | `source_type` | `string` | No | -- | Source type (e.g. `gps`, `router`, `bluetooth`, `bluetooth_le`) |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDeviceTracker
metadata:
  name: phone-tracker
spec:
  name: "Phone Tracker"
  stateTopic: "tracker/phone/state"
  sourceType: "bluetooth_le"
  payloadHome: "home"
  payloadNotHome: "not_home"
  device:
    name: "Phone Tracker"
    identifiers:
      - "phone-tracker-01"
    manufacturer: "Custom"
    model: "BLE Scanner"
```
