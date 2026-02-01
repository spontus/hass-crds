# MQTTClimate

A thermostat / HVAC entity with temperature control, modes, and fan speed.

- **Kind**: `MQTTClimate`
- **Resource**: `mqttclimates`
- **HA Component**: `climate`
- **HA Docs**: [MQTT Climate](https://www.home-assistant.io/integrations/climate.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `temperatureCommandTopic` | `temperature_command_topic` | `string` | No | -- | Topic to set target temperature |
| `temperatureStateTopic` | `temperature_state_topic` | `string` | No | -- | Topic to read target temperature |
| `temperatureCommandTemplate` | `temperature_command_template` | `string` | No | -- | Template for temperature command |
| `temperatureStateTemplate` | `temperature_state_template` | `string` | No | -- | Template to extract target temp |
| `currentTemperatureTopic` | `current_temperature_topic` | `string` | No | -- | Topic to read current temperature |
| `currentTemperatureTemplate` | `current_temperature_template` | `string` | No | -- | Template to extract current temp |
| `modeCommandTopic` | `mode_command_topic` | `string` | No | -- | Topic to set HVAC mode |
| `modeStateTopic` | `mode_state_topic` | `string` | No | -- | Topic to read HVAC mode |
| `modeCommandTemplate` | `mode_command_template` | `string` | No | -- | Template for mode command |
| `modeStateTemplate` | `mode_state_template` | `string` | No | -- | Template to extract mode |
| `modes` | `modes` | `[]string` | No | `["auto","off","cool","heat","dry","fan_only"]` | Supported HVAC modes |
| `fanModeCommandTopic` | `fan_mode_command_topic` | `string` | No | -- | Topic to set fan mode |
| `fanModeStateTopic` | `fan_mode_state_topic` | `string` | No | -- | Topic to read fan mode |
| `fanModeCommandTemplate` | `fan_mode_command_template` | `string` | No | -- | Template for fan mode command |
| `fanModeStateTemplate` | `fan_mode_state_template` | `string` | No | -- | Template to extract fan mode |
| `fanModes` | `fan_modes` | `[]string` | No | `["auto","low","medium","high"]` | Supported fan modes |
| `swingModeCommandTopic` | `swing_mode_command_topic` | `string` | No | -- | Topic to set swing mode |
| `swingModeStateTopic` | `swing_mode_state_topic` | `string` | No | -- | Topic to read swing mode |
| `swingModes` | `swing_modes` | `[]string` | No | `["on","off"]` | Supported swing modes |
| `presetModeCommandTopic` | `preset_mode_command_topic` | `string` | No | -- | Topic to set preset mode |
| `presetModeStateTopic` | `preset_mode_state_topic` | `string` | No | -- | Topic to read preset mode |
| `presetModes` | `preset_modes` | `[]string` | No | -- | Supported preset modes (e.g. `away`, `eco`, `boost`) |
| `actionTopic` | `action_topic` | `string` | No | -- | Topic to read current HVAC action |
| `actionTemplate` | `action_template` | `string` | No | -- | Template to extract action |
| `tempStep` | `temp_step` | `number` | No | `1` | Step size for temperature adjustments |
| `minTemp` | `min_temp` | `number` | No | -- | Minimum setpoint temperature |
| `maxTemp` | `max_temp` | `number` | No | -- | Maximum setpoint temperature |
| `temperatureUnit` | `temperature_unit` | `string` | No | -- | `C` or `F` |
| `precision` | `precision` | `number` | No | `0.1` | Temperature precision |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTClimate
metadata:
  name: living-room-thermostat
spec:
  name: "Living Room Thermostat"
  temperatureCommandTopic: "hvac/living-room/temp/set"
  temperatureStateTopic: "hvac/living-room/temp/state"
  currentTemperatureTopic: "hvac/living-room/temp/current"
  modeCommandTopic: "hvac/living-room/mode/set"
  modeStateTopic: "hvac/living-room/mode/state"
  modes:
    - "off"
    - "heat"
    - "cool"
    - "auto"
  minTemp: 16
  maxTemp: 30
  tempStep: 0.5
  device:
    name: "Living Room HVAC"
    identifiers:
      - "hvac-living-room-01"
```
