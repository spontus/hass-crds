# MQTTHumidifier

A humidifier entity with target humidity and mode support.

- **Kind**: `MQTTHumidifier`
- **Resource**: `mqtthumidifiers`
- **HA Component**: `humidifier`
- **HA Docs**: [MQTT Humidifier](https://www.home-assistant.io/integrations/humidifier.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish on/off commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current on/off state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload for on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload for off |
| `targetHumidityCommandTopic` | `target_humidity_command_topic` | `string` | Yes | -- | Topic to set target humidity |
| `targetHumidityStateTopic` | `target_humidity_state_topic` | `string` | No | -- | Topic to read target humidity |
| `targetHumidityCommandTemplate` | `target_humidity_command_template` | `string` | No | -- | Template for target humidity command |
| `targetHumidityStateTemplate` | `target_humidity_state_template` | `string` | No | -- | Template to extract target humidity |
| `currentHumidityTopic` | `current_humidity_topic` | `string` | No | -- | Topic to read current humidity |
| `currentHumidityTemplate` | `current_humidity_template` | `string` | No | -- | Template to extract current humidity |
| `modeCommandTopic` | `mode_command_topic` | `string` | No | -- | Topic to set mode |
| `modeStateTopic` | `mode_state_topic` | `string` | No | -- | Topic to read current mode |
| `modeCommandTemplate` | `mode_command_template` | `string` | No | -- | Template for mode command |
| `modeStateTemplate` | `mode_state_template` | `string` | No | -- | Template to extract mode |
| `modes` | `modes` | `[]string` | No | -- | Supported modes (e.g. `normal`, `eco`, `boost`, `sleep`) |
| `actionTopic` | `action_topic` | `string` | No | -- | Topic to read current action |
| `actionTemplate` | `action_template` | `string` | No | -- | Template to extract action |
| `minHumidity` | `min_humidity` | `number` | No | `0` | Minimum target humidity |
| `maxHumidity` | `max_humidity` | `number` | No | `100` | Maximum target humidity |
| `deviceClass` | `device_class` | `string` | No | -- | `humidifier` or `dehumidifier` |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTHumidifier
metadata:
  name: room-humidifier
spec:
  name: "Room Humidifier"
  commandTopic: "humidifier/room/set"
  stateTopic: "humidifier/room/state"
  targetHumidityCommandTopic: "humidifier/room/humidity/set"
  targetHumidityStateTopic: "humidifier/room/humidity/state"
  modes:
    - "normal"
    - "eco"
    - "boost"
    - "sleep"
  modeCommandTopic: "humidifier/room/mode/set"
  modeStateTopic: "humidifier/room/mode/state"
  minHumidity: 30
  maxHumidity: 80
  deviceClass: "humidifier"
  device:
    name: "Room Humidifier"
    identifiers:
      - "humidifier-room-01"
    manufacturer: "Custom"
    model: "Smart Humidifier"
```
