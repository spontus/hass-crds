# MQTTText

A free-text input entity.

- **Kind**: `MQTTText`
- **Resource**: `mqtttexts`
- **HA Component**: `text`
- **HA Docs**: [MQTT Text](https://www.home-assistant.io/integrations/text.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish text value |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current value |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract value from payload |
| `min` | `min` | `integer` | No | `0` | Minimum text length |
| `max` | `max` | `integer` | No | `255` | Maximum text length |
| `pattern` | `pattern` | `string` | No | -- | Regex pattern for validation |
| `mode` | `mode` | `string` | No | `text` | `text` or `password` |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTText
metadata:
  name: display-message
spec:
  name: "Display Message"
  commandTopic: "cmnd/display/message"
  stateTopic: "stat/display/message"
  max: 32
  icon: "mdi:message-text"
  device:
    name: "Dashboard Display"
    identifiers:
      - "dashboard-display-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTNumber](number.md), [MQTTSelect](select.md)
