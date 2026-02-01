# MQTTSensor

A read-only sensor that reports a value from an MQTT topic.

- **Kind**: `MQTTSensor`
- **Resource**: `mqttsensors`
- **HA Component**: `sensor`
- **HA Docs**: [MQTT Sensor](https://www.home-assistant.io/integrations/sensor.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `stateTopic` | `state_topic` | `string` | Yes | -- | Topic to read sensor value |
| `valueTemplate` | `value_template` | `string` | No | -- | Template to extract value from payload |
| `unitOfMeasurement` | `unit_of_measurement` | `string` | No | -- | Unit displayed in HA (e.g. `°C`, `%`, `W`) |
| `deviceClass` | `device_class` | `string` | No | -- | HA device class (e.g. `temperature`, `humidity`, `power`, `energy`, `battery`) |
| `stateClass` | `state_class` | `string` | No | -- | `measurement`, `total`, or `total_increasing` |
| `expireAfter` | `expire_after` | `integer` | No | -- | Seconds after which the sensor value expires |
| `forceUpdate` | `force_update` | `bool` | No | `false` | Update HA state even if the value hasn't changed |
| `lastResetValueTemplate` | `last_reset_value_template` | `string` | No | -- | Template for the last reset timestamp |
| `suggestedDisplayPrecision` | `suggested_display_precision` | `integer` | No | -- | Number of decimal places to display |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: cpu-temperature
spec:
  name: "CPU Temperature"
  stateTopic: "sensors/server/cpu_temp"
  unitOfMeasurement: "°C"
  deviceClass: "temperature"
  stateClass: "measurement"
  valueTemplate: "{{ value_json.temperature }}"
  suggestedDisplayPrecision: 1
  device:
    name: "Home Server"
    identifiers:
      - "home-server-01"
```
