# MQTTLock

A lock entity with optional code support.

- **Kind**: `MQTTLock`
- **Resource**: `mqttlocks`
- **HA Component**: `lock`
- **HA Docs**: [MQTT Lock](https://www.home-assistant.io/integrations/lock.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish lock/unlock commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current lock state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadLock` | `payload_lock` | `string` | No | `LOCK` | Payload for lock command |
| `payloadUnlock` | `payload_unlock` | `string` | No | `UNLOCK` | Payload for unlock command |
| `payloadOpen` | `payload_open` | `string` | No | -- | Payload for open command (unlatch) |
| `stateLocked` | `state_locked` | `string` | No | `LOCKED` | State value meaning locked |
| `stateUnlocked` | `state_unllocked` | `string` | No | `UNLOCKED` | State value meaning unlocked |
| `stateLocking` | `state_locking` | `string` | No | `LOCKING` | State value meaning locking |
| `stateUnlocking` | `state_unlocking` | `string` | No | `UNLOCKING` | State value meaning unlocking |
| `stateJammed` | `state_jammed` | `string` | No | `JAMMED` | State value meaning jammed |
| `codeFormat` | `code_format` | `string` | No | -- | Regex for valid codes (e.g. `^\d{4}$`) |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLock
metadata:
  name: front-door-lock
spec:
  name: "Front Door Lock"
  commandTopic: "locks/front-door/set"
  stateTopic: "locks/front-door/state"
  codeFormat: "^\\d{4}$"
  device:
    name: "Front Door Lock"
    identifiers:
      - "front-door-lock-01"
    manufacturer: "Custom"
    model: "Smart Lock v2"
```
