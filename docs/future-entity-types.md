# Entity Type Coverage

All 28 Home Assistant MQTT entity types are now supported.

## Implemented Entity Types

### Initial Release (10 types)

| HA Component | Kind | Docs |
|---|---|---|
| `button` | `MQTTButton` | [CRD](crds/button.md) |
| `switch` | `MQTTSwitch` | [CRD](crds/switch.md) |
| `light` | `MQTTLight` | [CRD](crds/light.md) |
| `sensor` | `MQTTSensor` | [CRD](crds/sensor.md) |
| `binary_sensor` | `MQTTBinarySensor` | [CRD](crds/binary-sensor.md) |
| `cover` | `MQTTCover` | [CRD](crds/cover.md) |
| `climate` | `MQTTClimate` | [CRD](crds/climate.md) |
| `number` | `MQTTNumber` | [CRD](crds/number.md) |
| `select` | `MQTTSelect` | [CRD](crds/select.md) |
| `text` | `MQTTText` | [CRD](crds/text.md) |

### Extended Release (18 types)

| HA Component | Kind | Docs |
|---|---|---|
| `fan` | `MQTTFan` | [CRD](crds/fan.md) |
| `lock` | `MQTTLock` | [CRD](crds/lock.md) |
| `siren` | `MQTTSiren` | [CRD](crds/siren.md) |
| `alarm_control_panel` | `MQTTAlarmControlPanel` | [CRD](crds/alarm-control-panel.md) |
| `humidifier` | `MQTTHumidifier` | [CRD](crds/humidifier.md) |
| `vacuum` | `MQTTVacuum` | [CRD](crds/vacuum.md) |
| `valve` | `MQTTValve` | [CRD](crds/valve.md) |
| `water_heater` | `MQTTWaterHeater` | [CRD](crds/water-heater.md) |
| `lawn_mower` | `MQTTLawnMower` | [CRD](crds/lawn-mower.md) |
| `camera` | `MQTTCamera` | [CRD](crds/camera.md) |
| `device_tracker` | `MQTTDeviceTracker` | [CRD](crds/device-tracker.md) |
| `device_trigger` | `MQTTDeviceTrigger` | [CRD](crds/device-trigger.md) |
| `event` | `MQTTEvent` | [CRD](crds/event.md) |
| `image` | `MQTTImage` | [CRD](crds/image.md) |
| `scene` | `MQTTScene` | [CRD](crds/scene.md) |
| `tag` | `MQTTTag` | [CRD](crds/tag.md) |
| `update` | `MQTTUpdate` | [CRD](crds/update.md) |
| `notify` | `MQTTNotify` | [CRD](crds/notify.md) |

## Contributing

When adding a new entity type:

1. Create the CRD definition in `config/crd/bases/`
2. Add the controller reconciliation logic
3. Add documentation in `docs/crds/<type>.md` following the existing pattern
4. Add an example in `docs/examples/`
5. Update `docs/crds/README.md` with the new type
6. Update RBAC rules in the controller deployment
