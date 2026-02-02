# CRD Reference

All CRDs belong to the API group `mqtt.home-assistant.io/v1alpha1`.

## Supported Entity Types

| Kind | Resource | HA Component | Description |
|---|---|---|---|
| [`MQTTButton`](button.md) | `mqttbuttons` | `button` | Stateless trigger |
| [`MQTTSwitch`](switch.md) | `mqttswitches` | `switch` | On/off toggle with state |
| [`MQTTLight`](light.md) | `mqttlights` | `light` | Light with optional brightness, color, effects |
| [`MQTTSensor`](sensor.md) | `mqttsensors` | `sensor` | Read-only sensor value |
| [`MQTTBinarySensor`](binary-sensor.md) | `mqttbinarysensors` | `binary_sensor` | On/off read-only sensor |
| [`MQTTCover`](cover.md) | `mqttcovers` | `cover` | Garage door, blind, shutter |
| [`MQTTClimate`](climate.md) | `mqttclimates` | `climate` | Thermostat / HVAC |
| [`MQTTNumber`](number.md) | `mqttnumbers` | `number` | Numeric input with min/max |
| [`MQTTSelect`](select.md) | `mqttselects` | `select` | Dropdown selection |
| [`MQTTText`](text.md) | `mqtttexts` | `text` | Free-text input |
| [`MQTTFan`](fan.md) | `mqttfans` | `fan` | Fan with speed, direction, oscillation |
| [`MQTTLock`](lock.md) | `mqttlocks` | `lock` | Lock/unlock with optional code |
| [`MQTTSiren`](siren.md) | `mqttsirens` | `siren` | Siren with tone and volume |
| [`MQTTAlarmControlPanel`](alarm-control-panel.md) | `mqttalarmcontrolpanels` | `alarm_control_panel` | Alarm panel with arm/disarm modes |
| [`MQTTHumidifier`](humidifier.md) | `mqtthumidifiers` | `humidifier` | Humidifier with target humidity and modes |
| [`MQTTVacuum`](vacuum.md) | `mqttvacuums` | `vacuum` | Robot vacuum with start/stop/return |
| [`MQTTValve`](valve.md) | `mqttvalves` | `valve` | Valve open/close (irrigation, plumbing) |
| [`MQTTWaterHeater`](water-heater.md) | `mqttwaterheaters` | `water_heater` | Water heater with temperature and modes |
| [`MQTTLawnMower`](lawn-mower.md) | `mqttlawnmowers` | `lawn_mower` | Robot lawn mower with start/dock |
| [`MQTTCamera`](camera.md) | `mqttcameras` | `camera` | Camera with image topic |
| [`MQTTDeviceTracker`](device-tracker.md) | `mqttdevicetrackers` | `device_tracker` | Device presence/location tracking |
| [`MQTTDeviceTrigger`](device-trigger.md) | `mqttdevicetriggers` | `device_trigger` | Device-based automation triggers |
| [`MQTTEvent`](event.md) | `mqttevents` | `event` | Event entity for stateless events |
| [`MQTTImage`](image.md) | `mqttimages` | `image` | Static image entity |
| [`MQTTScene`](scene.md) | `mqttscenes` | `scene` | Scene activation |
| [`MQTTTag`](tag.md) | `mqtttags` | `tag` | Tag scanner (NFC, RFID) |
| [`MQTTUpdate`](update.md) | `mqttupdates` | `update` | Firmware/software update entity |
| [`MQTTNotify`](notify.md) | `mqttnotifys` | `notify` | Notification service |

## Utility CRDs

| Kind | Resource | Description |
|---|---|---|
| [`MQTTDevice`](device.md) | `mqttdevices` | Shared device definition referenced by entity CRDs via `deviceRef` |

## Common Fields

All CRD types share a set of common fields for entity metadata, device association, availability, and MQTT options. These are documented in [Common Fields](common-fields.md).

## Field Naming Convention

CRD specs use **camelCase** (Kubernetes convention). The controller maps these to **snake_case** (Home Assistant MQTT convention) when building discovery payloads:

| CRD Spec | MQTT JSON |
|---|---|
| `commandTopic` | `command_topic` |
| `stateTopic` | `state_topic` |
| `uniqueId` | `unique_id` |
| `deviceClass` | `device_class` |
| `valueTemplate` | `value_template` |

## Discovery Topic

Each CRD instance publishes its discovery payload to:

```
<discovery_prefix>/<component>/<namespace>-<name>/config
```

For example, an `MQTTSensor` named `cpu-temp` in namespace `monitoring` publishes to:

```
homeassistant/sensor/monitoring-cpu-temp/config
```
