# MQTTTag

A tag scanner entity for NFC, RFID, or QR code scanning via MQTT.

- **Kind**: `MQTTTag`
- **Resource**: `mqtttags`
- **HA Component**: `tag`
- **HA Docs**: [MQTT Tag](https://www.home-assistant.io/integrations/tag.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `topic` | `topic` | `string` | Yes | -- | Topic to subscribe to for tag scans |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract tag ID from payload |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTTag
metadata:
  name: nfc-tag
spec:
  topic: "tags/nfc/scan"
  valueTemplate: "{{ value_json.tag_id }}"
  device:
    name: "NFC Reader"
    identifiers:
      - "nfc-reader-01"
    manufacturer: "Custom"
    model: "NFC Scanner"
    suggestedArea: "Entrance"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
- **Related**: [MQTTDeviceTracker](device-tracker.md), [MQTTEvent](event.md), [MQTTDeviceTrigger](device-trigger.md)
