# MQTTLight

A light entity with optional brightness, color temperature, and RGB color support. Supports three schema modes: basic (default), JSON, and template.

- **Kind**: `MQTTLight`
- **Resource**: `mqttlights`
- **HA Component**: `light`
- **HA Docs**: [MQTT Light](https://www.home-assistant.io/integrations/light.mqtt/)

## Schema Selection

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `schema` | `schema` | `string` | No | `default` | `default`, `json`, or `template` |

## Basic Schema Fields (`schema: default`)

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish on/off commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read current state |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload for on |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload for off |
| `brightnessCommandTopic` | `brightness_command_topic` | `string` | No | -- | Topic for brightness commands |
| `brightnessStateTopic` | `brightness_state_topic` | `string` | No | -- | Topic for brightness state |
| `brightnessScale` | `brightness_scale` | `integer` | No | `255` | Max brightness value |
| `brightnessValueTemplate` | `brightness_value_template` | `string` | No | -- | Template to extract brightness |
| `colorTempCommandTopic` | `color_temp_command_topic` | `string` | No | -- | Topic for color temperature commands |
| `colorTempStateTopic` | `color_temp_state_topic` | `string` | No | -- | Topic for color temperature state |
| `colorTempValueTemplate` | `color_temp_value_template` | `string` | No | -- | Template to extract color temp |
| `rgbCommandTopic` | `rgb_command_topic` | `string` | No | -- | Topic for RGB color commands |
| `rgbStateTopic` | `rgb_state_topic` | `string` | No | -- | Topic for RGB color state |
| `rgbCommandTemplate` | `rgb_command_template` | `string` | No | -- | Template for RGB command payload |
| `rgbValueTemplate` | `rgb_value_template` | `string` | No | -- | Template to extract RGB state |
| `effectCommandTopic` | `effect_command_topic` | `string` | No | -- | Topic for effect commands |
| `effectStateTopic` | `effect_state_topic` | `string` | No | -- | Topic for effect state |
| `effectList` | `effect_list` | `[]string` | No | -- | List of supported effects |
| `effectValueTemplate` | `effect_value_template` | `string` | No | -- | Template to extract effect |
| `minMireds` | `min_mireds` | `integer` | No | -- | Minimum color temp in mireds |
| `maxMireds` | `max_mireds` | `integer` | No | -- | Maximum color temp in mireds |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |
| `onCommandType` | `on_command_type` | `string` | No | -- | `last`, `first`, or `brightness` |

## JSON Schema Fields (`schema: json`)

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish JSON commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read JSON state |
| `brightness` | `brightness` | `bool` | No | `false` | Enable brightness support |
| `colorTemp` | `color_temp` | `bool` | No | `false` | Enable color temperature |
| `effect` | `effect` | `bool` | No | `false` | Enable effects |
| `effectList` | `effect_list` | `[]string` | No | -- | List of supported effects |
| `supportedColorModes` | `supported_color_modes` | `[]string` | No | -- | Supported color modes (e.g. `rgb`, `xy`, `hs`, `color_temp`) |
| `minMireds` | `min_mireds` | `integer` | No | -- | Minimum color temp in mireds |
| `maxMireds` | `max_mireds` | `integer` | No | -- | Maximum color temp in mireds |

## Template Schema Fields (`schema: template`)

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish commands |
| `stateTopic` | `state_topic` | `string` | No | -- | Topic to read state |
| `commandOnTemplate` | `command_on_template` | `string` | Yes | -- | Template for on command |
| `commandOffTemplate` | `command_off_template` | `string` | Yes | -- | Template for off command |
| `stateTemplate` | `state_template` | `string` | No | -- | Template to extract state |
| `brightnessTemplate` | `brightness_template` | `string` | No | -- | Template to extract brightness |
| `colorTempTemplate` | `color_temp_template` | `string` | No | -- | Template to extract color temp |
| `redTemplate` | `red_template` | `string` | No | -- | Template to extract red value |
| `greenTemplate` | `green_template` | `string` | No | -- | Template to extract green value |
| `blueTemplate` | `blue_template` | `string` | No | -- | Template to extract blue value |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example (JSON Schema with RGB)

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLight
metadata:
  name: led-strip
spec:
  name: "LED Strip"
  schema: "json"
  commandTopic: "lights/led-strip/set"
  stateTopic: "lights/led-strip/state"
  brightness: true
  supportedColorModes:
    - "rgb"
  effectList:
    - "rainbow"
    - "breathing"
    - "solid"
  device:
    name: "LED Controller"
    identifiers:
      - "led-controller-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTButton](button.md), [MQTTSwitch](switch.md), [MQTTFan](fan.md), [MQTTCover](cover.md), [MQTTLock](lock.md), [MQTTValve](valve.md)
