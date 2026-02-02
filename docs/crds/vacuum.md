# MQTTVacuum

A robot vacuum entity with start, stop, pause, return to base, and cleaning features.

- **Kind**: `MQTTVacuum`
- **Resource**: `mqttvacuums`
- **HA Component**: `vacuum`
- **HA Docs**: [MQTT Vacuum](https://www.home-assistant.io/integrations/vacuum.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | No | -- | Topic for basic commands (start, stop, return_to_base, etc.) |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read vacuum state |
| `sendCommandTopic` | `send_command_topic` | `string` | No | -- | Topic for custom commands |
| `setFanSpeedTopic` | `set_fan_speed_topic` | `string` | No | -- | Topic for fan speed commands |
| `fanSpeedList` | `fan_speed_list` | `[]string` | No | -- | List of supported fan speeds |
| `payloadStart` | `payload_start` | `string` | No | `start` | Payload for start command |
| `payloadStop` | `payload_stop` | `string` | No | `stop` | Payload for stop command |
| `payloadPause` | `payload_pause` | `string` | No | `pause` | Payload for pause command |
| `payloadReturnToBase` | `payload_return_to_base` | `string` | No | `return_to_base` | Payload for return to base command |
| `payloadCleanSpot` | `payload_clean_spot` | `string` | No | `clean_spot` | Payload for clean spot command |
| `payloadLocate` | `payload_locate` | `string` | No | `locate` | Payload for locate command |
| `supportedFeatures` | `supported_features` | `[]string` | No | -- | Supported features (e.g. `start`, `stop`, `pause`, `return_home`, `fan_speed`, `send_command`, `locate`, `clean_spot`) |
| `schema` | `schema` | `string` | No | `legacy` | `legacy` or `state` |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTVacuum
metadata:
  name: robot-vacuum
spec:
  name: "Robot Vacuum"
  schema: "state"
  commandTopic: "vacuum/robot/command"
  stateTopic: "vacuum/robot/state"
  setFanSpeedTopic: "vacuum/robot/fan_speed/set"
  sendCommandTopic: "vacuum/robot/custom_command"
  fanSpeedList:
    - "quiet"
    - "standard"
    - "turbo"
    - "max"
  supportedFeatures:
    - "start"
    - "stop"
    - "pause"
    - "return_home"
    - "fan_speed"
    - "send_command"
    - "locate"
  device:
    name: "Robot Vacuum"
    identifiers:
      - "robot-vacuum-01"
    manufacturer: "Custom"
    model: "CleanBot v3"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTLawnMower](lawn-mower.md), [MQTTUpdate](update.md)
