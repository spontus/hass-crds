# MQTTDevice

A shared device definition that entity CRDs can reference instead of duplicating the device block inline. `MQTTDevice` does not create an entity in Home Assistant -- it is a convenience resource that centralises device metadata for multiple entities.

- **Kind**: `MQTTDevice`
- **Resource**: `mqttdevices`
- **HA Component**: -- (no HA entity created)
- **HA Docs**: [Device Registry](https://developers.home-assistant.io/docs/device_registry_index/)

## Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `name` | `device.name` | `string` | No | -- | Device display name |
| `identifiers` | `device.identifiers` | `[]string` | No | -- | List of identifiers (at least one of `identifiers` or `connections` is needed) |
| `connections` | `device.connections` | `[][]string` | No | -- | List of `[type, value]` pairs (e.g. `[["mac", "aa:bb:cc:dd:ee:ff"]]`) |
| `manufacturer` | `device.manufacturer` | `string` | No | -- | Device manufacturer |
| `model` | `device.model` | `string` | No | -- | Device model |
| `modelId` | `device.model_id` | `string` | No | -- | Device model identifier |
| `serialNumber` | `device.serial_number` | `string` | No | -- | Device serial number |
| `hwVersion` | `device.hw_version` | `string` | No | -- | Hardware version |
| `swVersion` | `device.sw_version` | `string` | No | -- | Software version |
| `suggestedArea` | `device.suggested_area` | `string` | No | -- | Suggested area in HA (e.g. "Living Room") |
| `configurationUrl` | `device.configuration_url` | `string` | No | -- | URL for device configuration |
| `viaDevice` | `device.via_device` | `string` | No | -- | Identifier of device that routes messages |

## Usage

Entity CRDs reference an `MQTTDevice` by name using `deviceRef` instead of an inline `device` block:

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDevice
metadata:
  name: weather-station
  namespace: hass-crds
spec:
  name: "Weather Station"
  identifiers:
    - "weather-station-01"
  manufacturer: "Custom"
  model: "ESP32 Weather Station"
  swVersion: "2.1.0"
  suggestedArea: "Garden"
---
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: weather-temperature
  namespace: hass-crds
spec:
  name: "Temperature"
  stateTopic: "sensors/weather/temperature"
  unitOfMeasurement: "°C"
  deviceClass: "temperature"
  deviceRef:
    name: "weather-station"
---
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: weather-humidity
  namespace: hass-crds
spec:
  name: "Humidity"
  stateTopic: "sensors/weather/humidity"
  unitOfMeasurement: "%"
  deviceClass: "humidity"
  deviceRef:
    name: "weather-station"
```

When the controller encounters `deviceRef`, it resolves the referenced `MQTTDevice` and injects its fields into the discovery payload's `device` block. If both `device` and `deviceRef` are set, the inline `device` block takes priority.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDevice
metadata:
  name: sensor-hub
  namespace: hass-crds
spec:
  name: "Living Room Sensor Hub"
  identifiers:
    - "sensor-hub-01"
  manufacturer: "Custom"
  model: "SensorHub v2"
  swVersion: "1.4.0"
  suggestedArea: "Living Room"
  configurationUrl: "http://sensor-hub-01.local"
```

---

## See Also

- [CRD Reference](README.md) - All entity types
- [Common Fields](common-fields.md) - Shared fields (device, availability, MQTT options)
