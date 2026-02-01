# MQTTSwitch

An on/off toggle entity with state feedback.

- **Kind**: `MQTTSwitch`
- **Resource**: `mqttswitch`
- **HA Component**: `switch`
- **HA Docs**: [MQTT Switch](https://www.home-assistant.io/integrations/switch.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish on/off commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload representing on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload representing off |
| `stateOn` | `state_on` | `string` | No | -- | State value that means on (if different from payloadOn) |
| `stateOff` | `state_off` | `string` | No | -- | State value that means off (if different from payloadOff) |
| `deviceClass` | `device_class` | `string` | No | -- | `outlet` or `switch` |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSwitch
metadata:
  name: desk-lamp
spec:
  name: "Desk Lamp"
  commandTopic: "cmnd/desk-lamp/POWER"
  stateTopic: "stat/desk-lamp/POWER"
  payloadOn: "ON"
  payloadOff: "OFF"
  device:
    name: "Desk Lamp"
    identifiers:
      - "desk-lamp-01"
    manufacturer: "Sonoff"
    model: "Basic R2"
```
