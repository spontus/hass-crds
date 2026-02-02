# MQTTNumber

A numeric input entity with min/max bounds and step size.

- **Kind**: `MQTTNumber`
- **Resource**: `mqttnumbers`
- **HA Component**: `number`
- **HA Docs**: [MQTT Number](https://www.home-assistant.io/integrations/number.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish number value |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current value |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract value from payload |
| `min` | `min` | `number` | No | `1` | Minimum value |
| `max` | `max` | `number` | No | `100` | Maximum value |
| `step` | `step` | `number` | No | `1` | Step size |
| `mode` | `mode` | `string` | No | `auto` | UI mode: `auto`, `box`, or `slider` |
| `unitOfMeasurement` | `unit_of_measurement` | `string` | No | -- | Unit displayed in HA |
| `deviceClass` | `device_class` | `string` | No | -- | HA device class (e.g. `temperature`, `humidity`, `power_factor`) |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTNumber
metadata:
  name: display-brightness
spec:
  name: "Display Brightness"
  commandTopic: "cmnd/display/brightness"
  stateTopic: "stat/display/brightness"
  min: 0
  max: 100
  step: 5
  mode: "slider"
  unitOfMeasurement: "%"
  icon: "mdi:brightness-6"
  device:
    name: "Dashboard Display"
    identifiers:
      - "dashboard-display-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTSelect](select.md), [MQTTText](text.md)
