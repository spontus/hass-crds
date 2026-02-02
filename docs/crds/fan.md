# MQTTFan

A fan entity with speed, direction, oscillation, and preset mode support.

- **Kind**: `MQTTFan`
- **Resource**: `mqttfans`
- **HA Component**: `fan`
- **HA Docs**: [MQTT Fan](https://www.home-assistant.io/integrations/fan.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish on/off commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current on/off state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload for on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload for off |
| `percentageCommandTopic` | `percentage_command_topic` | `string` | No | -- | Topic for speed percentage commands |
| `percentageStateTopic` | `percentage_state_topic` | `string` | No | -- | Topic to read speed percentage |
| `percentageCommandTemplate` | `percentage_command_template` | `string` | No | -- | Template for percentage command |
| `percentageValueTemplate` | `percentage_value_template` | `string` | No | -- | Template to extract percentage |
| `speedRangeMin` | `speed_range_min` | `integer` | No | `1` | Minimum speed value |
| `speedRangeMax` | `speed_range_max` | `integer` | No | `100` | Maximum speed value |
| `presetModeCommandTopic` | `preset_mode_command_topic` | `string` | No | -- | Topic for preset mode commands |
| `presetModeStateTopic` | `preset_mode_state_topic` | `string` | No | -- | Topic to read preset mode |
| `presetModeCommandTemplate` | `preset_mode_command_template` | `string` | No | -- | Template for preset mode command |
| `presetModeValueTemplate` | `preset_mode_value_template` | `string` | No | -- | Template to extract preset mode |
| `presetModes` | `preset_modes` | `[]string` | No | -- | List of supported preset modes |
| `oscillationCommandTopic` | `oscillation_command_topic` | `string` | No | -- | Topic for oscillation commands |
| `oscillationStateTopic` | `oscillation_state_topic` | `string` | No | -- | Topic to read oscillation state |
| `oscillationCommandTemplate` | `oscillation_command_template` | `string` | No | -- | Template for oscillation command |
| `oscillationValueTemplate` | `oscillation_value_template` | `string` | No | -- | Template to extract oscillation state |
| `payloadOscillationOn` | `payload_oscillation_on` | `string` | No | `oscillate_on` | Payload for oscillation on |
| `payloadOscillationOff` | `payload_oscillation_off` | `string` | No | `oscillate_off` | Payload for oscillation off |
| `directionCommandTopic` | `direction_command_topic` | `string` | No | -- | Topic for direction commands |
| `directionStateTopic` | `direction_state_topic` | `string` | No | -- | Topic to read direction state |
| `directionValueTemplate` | `direction_value_template` | `string` | No | -- | Template to extract direction |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTFan
metadata:
  name: ceiling-fan
spec:
  name: "Ceiling Fan"
  commandTopic: "fans/ceiling/set"
  stateTopic: "fans/ceiling/state"
  percentageCommandTopic: "fans/ceiling/speed/set"
  percentageStateTopic: "fans/ceiling/speed/state"
  speedRangeMin: 1
  speedRangeMax: 6
  presetModes:
    - "normal"
    - "breeze"
    - "sleep"
  presetModeCommandTopic: "fans/ceiling/preset/set"
  presetModeStateTopic: "fans/ceiling/preset/state"
  oscillationCommandTopic: "fans/ceiling/oscillation/set"
  oscillationStateTopic: "fans/ceiling/oscillation/state"
  directionCommandTopic: "fans/ceiling/direction/set"
  directionStateTopic: "fans/ceiling/direction/state"
  device:
    name: "Ceiling Fan"
    identifiers:
      - "ceiling-fan-01"
    manufacturer: "Custom"
    model: "Smart Fan Controller"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTButton](button.md), [MQTTSwitch](switch.md), [MQTTLight](light.md), [MQTTCover](cover.md), [MQTTLock](lock.md), [MQTTValve](valve.md)
