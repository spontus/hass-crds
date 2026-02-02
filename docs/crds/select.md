# MQTTSelect

A dropdown selection entity with a fixed list of options.

- **Kind**: `MQTTSelect`
- **Resource**: `mqttselects`
- **HA Component**: `select`
- **HA Docs**: [MQTT Select](https://www.home-assistant.io/integrations/select.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish selected option |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current selection |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract value from payload |
| `options` | `options` | `[]string` | Yes | -- | List of selectable options |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSelect
metadata:
  name: fan-speed
spec:
  name: "Fan Speed"
  commandTopic: "cmnd/fan/speed"
  stateTopic: "stat/fan/speed"
  options:
    - "off"
    - "low"
    - "medium"
    - "high"
  icon: "mdi:fan"
  device:
    name: "Desk Fan"
    identifiers:
      - "desk-fan-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTNumber](number.md), [MQTTText](text.md)
