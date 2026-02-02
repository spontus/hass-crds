# MQTTDeviceTrigger

A device automation trigger that fires when a specific MQTT message is received. Unlike other entity types, device triggers do not create entities in Home Assistant but register as device automations.

- **Kind**: `MQTTDeviceTrigger`
- **Resource**: `mqttdevicetriggers`
- **HA Component**: `device_trigger`
- **HA Docs**: [MQTT Device Trigger](https://www.home-assistant.io/integrations/device_trigger.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `topic` | `topic` | `string` | Yes | -- | MQTT topic to subscribe to for trigger events |
| `payload` | `payload` | `string` | No | -- | Specific payload that triggers the automation |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract value from payload |
| `type` | `type` | `string` | Yes | -- | Trigger type (e.g. `button_short_press`, `button_long_press`) |
| `subtype` | `subtype` | `string` | Yes | -- | Trigger subtype (e.g. `button_1`, `turn_on`) |
| `automationType` | `automation_type` | `string` | No | `trigger` | Automation type (always `trigger`) |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDeviceTrigger
metadata:
  name: button-trigger
spec:
  topic: "zigbee/button/action"
  type: "button_short_press"
  subtype: "button_1"
  automationType: "trigger"
  device:
    name: "Zigbee Button"
    identifiers:
      - "zigbee-button-01"
    manufacturer: "Custom"
    model: "Wireless Button"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTDeviceTracker](device-tracker.md), [MQTTTag](tag.md), [MQTTEvent](event.md)
