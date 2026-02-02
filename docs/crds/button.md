# MQTTButton

A stateless button entity. When pressed in Home Assistant, it publishes a payload to the command topic.

- **Kind**: `MQTTButton`
- **Resource**: `mqttbuttons`
- **HA Component**: `button`
- **HA Docs**: [MQTT Button](https://www.home-assistant.io/integrations/button.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish when button is pressed |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the command payload |
| `payloadPress` | `payload_press` | `string` | No | `PRESS` | Payload sent when button is pressed |
| `deviceClass` | `device_class` | `string` | No | -- | `identify`, `restart`, or `update` |

In addition to the fields above, all [common fields](common-fields.md) (entity metadata, device, availability, MQTT options, JSON attributes) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: restart-server
spec:
  name: "Restart Server"
  commandTopic: "cmnd/server/restart"
  payloadPress: "RESTART"
  deviceClass: "restart"
  icon: "mdi:restart"
  device:
    name: "Home Server"
    identifiers:
      - "home-server-01"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTSwitch](switch.md), [MQTTLight](light.md), [MQTTFan](fan.md), [MQTTCover](cover.md), [MQTTLock](lock.md), [MQTTValve](valve.md)
