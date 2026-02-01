# MQTTValve

A valve entity for controlling water, gas, or irrigation valves.

- **Kind**: `MQTTValve`
- **Resource**: `mqttvalves`
- **HA Component**: `valve`
- **HA Docs**: [MQTT Valve](https://www.home-assistant.io/integrations/valve.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | No | -- | Topic to publish open/close commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current valve state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `positionTopic` | `position_topic` | `string` | No | -- | Topic to read current position |
| `setPositionTopic` | `set_position_topic` | `string` | No | -- | Topic to publish position commands |
| `setPositionTemplate` | `set_position_template` | `string` | No | -- | Template for position command payload |
| `positionTemplate` | `position_template` | `string` | No | -- | Template to extract position from payload |
| `payloadOpen` | `payload_open` | `string` | No | `OPEN` | Payload for open command |
| `payloadClose` | `payload_close` | `string` | No | `CLOSE` | Payload for close command |
| `payloadStop` | `payload_stop` | `string` | No | `STOP` | Payload for stop command |
| `stateOpen` | `state_open` | `string` | No | `open` | State value meaning open |
| `stateClosed` | `state_closed` | `string` | No | `closed` | State value meaning closed |
| `stateOpening` | `state_opening` | `string` | No | `opening` | State value meaning opening |
| `stateClosing` | `state_closing` | `string` | No | `closing` | State value meaning closing |
| `deviceClass` | `device_class` | `string` | No | -- | `water`, `gas` |
| `reportsPosition` | `reports_position` | `bool` | No | `false` | Whether the valve reports position |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTValve
metadata:
  name: irrigation-valve
spec:
  name: "Irrigation Valve"
  commandTopic: "valves/irrigation/set"
  stateTopic: "valves/irrigation/state"
  deviceClass: "water"
  device:
    name: "Garden Irrigation"
    identifiers:
      - "irrigation-valve-01"
    manufacturer: "Custom"
    model: "Smart Valve"
    suggestedArea: "Garden"
```
