# MQTTSiren

A siren entity with optional tone, volume, and duration support.

- **Kind**: `MQTTSiren`
- **Resource**: `mqttsirens`
- **HA Component**: `siren`
- **HA Docs**: [MQTT Siren](https://www.home-assistant.io/integrations/siren.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish on/off commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload for on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload for off |
| `stateOn` | `state_on` | `string` | No | `ON` | State value meaning on |
| `stateOff` | `state_off` | `string` | No | `OFF` | State value meaning off |
| `availableTones` | `available_tones` | `[]string` | No | -- | List of supported tones |
| `supportTurnOn` | `support_turn_on` | `bool` | No | `true` | Whether the siren supports turn on |
| `supportTurnOff` | `support_turn_off` | `bool` | No | `true` | Whether the siren supports turn off |
| `supportDuration` | `support_duration` | `bool` | No | `true` | Whether duration is supported |
| `supportVolumeSet` | `support_volume_set` | `bool` | No | `true` | Whether volume level is supported |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSiren
metadata:
  name: alarm-siren
spec:
  name: "Alarm Siren"
  commandTopic: "siren/alarm/set"
  stateTopic: "siren/alarm/state"
  availableTones:
    - "fire"
    - "intruder"
    - "doorbell"
  device:
    name: "Alarm Siren"
    identifiers:
      - "alarm-siren-01"
    manufacturer: "Custom"
    model: "Multi-Tone Siren"
```
