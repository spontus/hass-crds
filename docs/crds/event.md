# MQTTEvent

An event entity for stateless events such as button presses or doorbell rings.

- **Kind**: `MQTTEvent`
- **Resource**: `mqttevents`
- **HA Component**: `event`
- **HA Docs**: [MQTT Event](https://www.home-assistant.io/integrations/event.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic to subscribe to for events |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract event type from payload |
| `eventTypes` | `event_types` | `[]string` | Yes | -- | List of supported event types |
| `deviceClass` | `device_class` | `string` | No | -- | `button`, `doorbell`, `motion` |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTEvent
metadata:
  name: doorbell-event
spec:
  name: "Doorbell"
  stateTopic: "doorbell/event"
  eventTypes:
    - "press"
    - "double_press"
  deviceClass: "doorbell"
  device:
    name: "Smart Doorbell"
    identifiers:
      - "doorbell-01"
    manufacturer: "Custom"
    model: "Smart Doorbell v2"
    suggestedArea: "Front Door"
```
