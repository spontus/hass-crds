# MQTTCover

A cover entity for garage doors, blinds, shutters, and similar devices.

- **Kind**: `MQTTCover`
- **Resource**: `mqttcovers`
- **HA Component**: `cover`
- **HA Docs**: [MQTT Cover](https://www.home-assistant.io/integrations/cover.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | No | -- | Topic for open/close/stop commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read cover state |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `positionTopic` | `position_topic` | `string` | No | -- | Topic to read current position |
| `setPositionTopic` | `set_position_topic` | `string` | No | -- | Topic to publish position commands |
| `setPositionTemplate` | `set_position_template` | `string` | No | -- | Template for position command payload |
| `positionTemplate` | `position_template` | `string` | No | -- | Template to extract position from payload |
| `tiltCommandTopic` | `tilt_command_topic` | `string` | No | -- | Topic for tilt commands |
| `tiltStatusTopic` | `tilt_status_topic` | `string` | No | -- | Topic to read tilt position |
| `tiltStatusTemplate` | `tilt_status_template` | `string` | No | -- | Template to extract tilt from payload |
| `payloadOpen` | `payload_open` | `string` | No | `OPEN` | Payload for open command |
| `payloadClose` | `payload_close` | `string` | No | `CLOSE` | Payload for close command |
| `payloadStop` | `payload_stop` | `string` | No | `STOP` | Payload for stop command |
| `stateOpen` | `state_open` | `string` | No | `open` | State value meaning open |
| `stateClosed` | `state_closed` | `string` | No | `closed` | State value meaning closed |
| `stateOpening` | `state_opening` | `string` | No | `opening` | State value meaning opening |
| `stateClosing` | `state_closing` | `string` | No | `closing` | State value meaning closing |
| `stateStopped` | `state_stopped` | `string` | No | `stopped` | State value meaning stopped |
| `positionOpen` | `position_open` | `integer` | No | `100` | Position value for fully open |
| `positionClosed` | `position_closed` | `integer` | No | `0` | Position value for fully closed |
| `tiltMin` | `tilt_min` | `integer` | No | `0` | Minimum tilt value |
| `tiltMax` | `tilt_max` | `integer` | No | `100` | Maximum tilt value |
| `deviceClass` | `device_class` | `string` | No | -- | `awning`, `blind`, `curtain`, `damper`, `door`, `garage`, `gate`, `shade`, `shutter`, `window` |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTCover
metadata:
  name: garage-door
spec:
  name: "Garage Door"
  commandTopic: "cmnd/garage/door"
  stateTopic: "stat/garage/door"
  deviceClass: "garage"
  payloadOpen: "OPEN"
  payloadClose: "CLOSE"
  payloadStop: "STOP"
  device:
    name: "Garage Controller"
    identifiers:
      - "garage-controller-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTButton](button.md), [MQTTSwitch](switch.md), [MQTTLight](light.md), [MQTTFan](fan.md), [MQTTLock](lock.md), [MQTTValve](valve.md)
