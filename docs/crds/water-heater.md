# MQTTWaterHeater

A water heater entity with temperature control and operation modes.

- **Kind**: `MQTTWaterHeater`
- **Resource**: `mqttwaterheaters`
- **HA Component**: `water_heater`
- **HA Docs**: [MQTT Water Heater](https://www.home-assistant.io/integrations/water_heater.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `temperatureCommandTopic` | `temperature_command_topic` | `string` | No | -- | Topic to set target temperature |
| `temperatureStateTopic` | `temperature_state_topic` | `string` | No | -- | Topic to read target temperature |
| `temperatureCommandTemplate` | `temperature_command_template` | `string` | No | -- | Template for temperature command |
| `temperatureStateTemplate` | `temperature_state_template` | `string` | No | -- | Template to extract target temp |
| `currentTemperatureTopic` | `current_temperature_topic` | `string` | No | -- | Topic to read current temperature |
| `currentTemperatureTemplate` | `current_temperature_template` | `string` | No | -- | Template to extract current temp |
| `modeCommandTopic` | `mode_command_topic` | `string` | No | -- | Topic to set operation mode |
| `modeStateTopic` | `mode_state_topic` | `string` | No | -- | Topic to read operation mode |
| `modeCommandTemplate` | `mode_command_template` | `string` | No | -- | Template for mode command |
| `modeStateTemplate` | `mode_state_template` | `string` | No | -- | Template to extract mode |
| `modes` | `modes` | `[]string` | No | -- | Supported modes (e.g. `off`, `eco`, `electric`, `gas`, `heat_pump`, `high_demand`, `performance`) |
| `powerCommandTopic` | `power_command_topic` | `string` | No | -- | Topic to publish on/off commands |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload for on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload for off |
| `minTemp` | `min_temp` | `number` | No | `110` | Minimum target temperature |
| `maxTemp` | `max_temp` | `number` | No | `140` | Maximum target temperature |
| `temperatureUnit` | `temperature_unit` | `string` | No | -- | `C` or `F` |
| `precision` | `precision` | `number` | No | `0.1` | Temperature precision |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTWaterHeater
metadata:
  name: water-heater
spec:
  name: "Water Heater"
  temperatureCommandTopic: "water-heater/temp/set"
  temperatureStateTopic: "water-heater/temp/state"
  currentTemperatureTopic: "water-heater/temp/current"
  modeCommandTopic: "water-heater/mode/set"
  modeStateTopic: "water-heater/mode/state"
  modes:
    - "off"
    - "eco"
    - "electric"
    - "performance"
  minTemp: 30
  maxTemp: 70
  temperatureUnit: "C"
  device:
    name: "Water Heater"
    identifiers:
      - "water-heater-01"
    manufacturer: "Custom"
    model: "Smart Water Heater"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTClimate](climate.md), [MQTTHumidifier](humidifier.md)
