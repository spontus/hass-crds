# MQTTNotify

A notification service entity that sends messages to a device via MQTT.

- **Kind**: `MQTTNotify`
- **Resource**: `mqttnotifys`
- **HA Component**: `notify`
- **HA Docs**: [MQTT Notify](https://www.home-assistant.io/integrations/notify.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `commandTopic` | `command_topic` | `string` | Yes | -- | Topic to publish notification messages |
| `commandTemplate` | `command_template` | `string` | No | -- | Template for the notification payload |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTNotify
metadata:
  name: display-notify
spec:
  name: "Display Notification"
  commandTopic: "display/notify"
  commandTemplate: '{"text": "{{ value }}"}'
  device:
    name: "Smart Display"
    identifiers:
      - "smart-display-01"
    manufacturer: "Custom"
    model: "E-Ink Display"
    suggestedArea: "Kitchen"
```
