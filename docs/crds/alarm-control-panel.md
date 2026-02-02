# MQTTAlarmControlPanel

An alarm control panel entity with arm/disarm modes and optional code support.

- **Kind**: `MQTTAlarmControlPanel`
- **Resource**: `mqttalarmcontrolpanels`
- **HA Component**: `alarm_control_panel`
- **HA Docs**: [MQTT Alarm Control Panel](https://www.home-assistant.io/integrations/alarm_control_panel.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish arm/disarm commands |
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic to read alarm state |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadArmHome` | `payload_arm_home` | `string` | No | `ARM_HOME` | Payload for arm home |
| `payloadArmAway` | `payload_arm_away` | `string` | No | `ARM_AWAY` | Payload for arm away |
| `payloadArmNight` | `payload_arm_night` | `string` | No | `ARM_NIGHT` | Payload for arm night |
| `payloadArmVacation` | `payload_arm_vacation` | `string` | No | `ARM_VACATION` | Payload for arm vacation |
| `payloadArmCustomBypass` | `payload_arm_custom_bypass` | `string` | No | `ARM_CUSTOM_BYPASS` | Payload for arm custom bypass |
| `payloadDisarm` | `payload_disarm` | `string` | No | `DISARM` | Payload for disarm |
| `payloadTrigger` | `payload_trigger` | `string` | No | -- | Payload for trigger |
| `codeArmRequired` | `code_arm_required` | `bool` | No | `true` | Whether code is required to arm |
| `codeDisarmRequired` | `code_disarm_required` | `bool` | No | `true` | Whether code is required to disarm |
| `codeTriggerRequired` | `code_trigger_required` | `bool` | No | `true` | Whether code is required to trigger |
| `codeFormat` | `code_format` | `string` | No | -- | `number` or `text` |
| `supportedFeatures` | `supported_features` | `[]string` | No | -- | Supported features (e.g. `arm_home`, `arm_away`, `arm_night`, `trigger`) |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTAlarmControlPanel
metadata:
  name: home-alarm
spec:
  name: "Home Alarm"
  commandTopic: "alarm/home/set"
  stateTopic: "alarm/home/state"
  codeFormat: "number"
  codeArmRequired: true
  codeDisarmRequired: true
  supportedFeatures:
    - "arm_home"
    - "arm_away"
    - "arm_night"
    - "trigger"
  device:
    name: "Home Alarm Panel"
    identifiers:
      - "alarm-panel-01"
    manufacturer: "Custom"
    model: "Smart Alarm v1"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTSiren](siren.md), [MQTTNotify](notify.md)
