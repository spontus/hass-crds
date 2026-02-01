# Examples

Copy-pasteable YAML manifests for every supported entity type. All examples use the `hass-crds` namespace -- adjust as needed.

| Example | Entity Type | Description |
|---|---|---|
| [basic-button.yaml](basic-button.yaml) | MQTTButton | Simple restart button |
| [switch-with-state.yaml](switch-with-state.yaml) | MQTTSwitch | Switch with state feedback |
| [rgb-light.yaml](rgb-light.yaml) | MQTTLight | JSON-schema light with RGB and effects |
| [temperature-sensor.yaml](temperature-sensor.yaml) | MQTTSensor | Temperature sensor with device class |
| [motion-binary-sensor.yaml](motion-binary-sensor.yaml) | MQTTBinarySensor | Motion detector with auto-reset |
| [garage-cover.yaml](garage-cover.yaml) | MQTTCover | Garage door with position |
| [thermostat-climate.yaml](thermostat-climate.yaml) | MQTTClimate | Thermostat with heat/cool modes |
| [brightness-number.yaml](brightness-number.yaml) | MQTTNumber | Brightness slider |
| [input-select.yaml](input-select.yaml) | MQTTSelect | Mode selector dropdown |
| [device-grouping.yaml](device-grouping.yaml) | Multiple | Multiple entities sharing one device |
| [ceiling-fan.yaml](ceiling-fan.yaml) | MQTTFan | Fan with speed/direction/oscillation |
| [front-door-lock.yaml](front-door-lock.yaml) | MQTTLock | Lock with code |
| [alarm-siren.yaml](alarm-siren.yaml) | MQTTSiren | Siren with tones |
| [home-alarm.yaml](home-alarm.yaml) | MQTTAlarmControlPanel | Alarm panel |
| [room-humidifier.yaml](room-humidifier.yaml) | MQTTHumidifier | Humidifier with modes |
| [robot-vacuum.yaml](robot-vacuum.yaml) | MQTTVacuum | Vacuum with features |
| [irrigation-valve.yaml](irrigation-valve.yaml) | MQTTValve | Garden valve |
| [water-heater.yaml](water-heater.yaml) | MQTTWaterHeater | Water heater with modes |
| [robot-mower.yaml](robot-mower.yaml) | MQTTLawnMower | Lawn mower |
| [doorbell-camera.yaml](doorbell-camera.yaml) | MQTTCamera | Camera with image topic |
| [phone-tracker.yaml](phone-tracker.yaml) | MQTTDeviceTracker | Device tracker |
| [button-trigger.yaml](button-trigger.yaml) | MQTTDeviceTrigger | Button press trigger |
| [doorbell-event.yaml](doorbell-event.yaml) | MQTTEvent | Doorbell press event |
| [plant-image.yaml](plant-image.yaml) | MQTTImage | Plant cam image |
| [movie-night-scene.yaml](movie-night-scene.yaml) | MQTTScene | Scene activation |
| [nfc-tag.yaml](nfc-tag.yaml) | MQTTTag | NFC tag scanner |
| [firmware-update.yaml](firmware-update.yaml) | MQTTUpdate | Firmware update entity |
| [display-notify.yaml](display-notify.yaml) | MQTTNotify | Notification service |
| [shared-device.yaml](shared-device.yaml) | MQTTDevice | Shared device with multiple entity references |

## Usage

```bash
# Apply a single example
kubectl apply -f basic-button.yaml

# Apply all examples
kubectl apply -f .
```
