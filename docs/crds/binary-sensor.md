# MQTTBinarySensor

A read-only on/off sensor (e.g. motion detector, door contact, leak sensor).

- **Kind**: `MQTTBinarySensor`
- **Resource**: `mqttbinarysensors`
- **HA Component**: `binary_sensor`
- **HA Docs**: [MQTT Binary Sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic to read sensor state |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract state from payload |
| `payloadOn` | `payload_on` | `string` | No | `ON` | Payload representing on/detected |
| `payloadOff` | `payload_off` | `string` | No | `OFF` | Payload representing off/clear |
| `deviceClass` | `device_class` | `string` | No | -- | HA device class (e.g. `motion`, `door`, `window`, `moisture`, `smoke`, `occupancy`) |
| `expireAfter` | `expire_after` | `integer` | No | -- | Seconds after which the state expires |
| `forceUpdate` | `force_update` | `bool` | No | `false` | Update state even if unchanged |
| `offDelay` | `off_delay` | `integer` | No | -- | Seconds after which the sensor auto-resets to off |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTBinarySensor
metadata:
  name: front-door
spec:
  name: "Front Door"
  stateTopic: "sensors/front-door/contact"
  deviceClass: "door"
  payloadOn: "OPEN"
  payloadOff: "CLOSED"
  device:
    name: "Front Door Sensor"
    identifiers:
      - "front-door-sensor-01"
```
