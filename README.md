# hass-crds

A Kubernetes operator that manages Home Assistant entities via MQTT autodiscovery. Define your HA entities as Kubernetes Custom Resources, and the controller automatically publishes MQTT discovery payloads so Home Assistant creates them.

## Features

- **28 entity types supported**: Button, Sensor, Switch, Light, Climate, Cover, Lock, Fan, Vacuum, and more
- **Declarative configuration**: Define entities as Kubernetes CRDs
- **Automatic discovery**: Controller publishes MQTT discovery payloads
- **Clean deletion**: Entities are removed from Home Assistant when CRDs are deleted
- **Device grouping**: Group multiple entities under a single device
- **Availability tracking**: Configure availability topics for entity status

## Supported Entity Types

| Entity Type | CRD Kind |
|-------------|----------|
| Alarm Control Panel | `MQTTAlarmControlPanel` |
| Binary Sensor | `MQTTBinarySensor` |
| Button | `MQTTButton` |
| Camera | `MQTTCamera` |
| Climate | `MQTTClimate` |
| Cover | `MQTTCover` |
| Device Tracker | `MQTTDeviceTracker` |
| Device Trigger | `MQTTDeviceTrigger` |
| Event | `MQTTEvent` |
| Fan | `MQTTFan` |
| Humidifier | `MQTTHumidifier` |
| Image | `MQTTImage` |
| Lawn Mower | `MQTTLawnMower` |
| Light | `MQTTLight` |
| Lock | `MQTTLock` |
| Notify | `MQTTNotify` |
| Number | `MQTTNumber` |
| Scene | `MQTTScene` |
| Select | `MQTTSelect` |
| Sensor | `MQTTSensor` |
| Siren | `MQTTSiren` |
| Switch | `MQTTSwitch` |
| Tag | `MQTTTag` |
| Text | `MQTTText` |
| Update | `MQTTUpdate` |
| Vacuum | `MQTTVacuum` |
| Valve | `MQTTValve` |
| Water Heater | `MQTTWaterHeater` |

## Installation

### Prerequisites

- Kubernetes cluster (1.26+)
- MQTT broker accessible from the cluster (e.g., Mosquitto)
- Home Assistant configured with MQTT integration

### Quick Install

```bash
# Install CRDs
kubectl apply -f https://raw.githubusercontent.com/spontus/hass-crds/main/config/crd/crds.yaml

# Deploy controller (creates namespace hass-crds-system)
kubectl apply -f https://raw.githubusercontent.com/spontus/hass-crds/main/dist/install.yaml

# Create MQTT configuration secret
kubectl create secret generic mqtt-config -n hass-crds-system \
  --from-literal=MQTT_BROKER=mqtt.example.com \
  --from-literal=MQTT_PORT=1883 \
  --from-literal=MQTT_USERNAME=homeassistant \
  --from-literal=MQTT_PASSWORD=your-password

# Restart controller to pick up secret
kubectl rollout restart deployment/hass-crds-controller-manager -n hass-crds-system
```

### Install from Source

```bash
git clone https://github.com/spontus/hass-crds.git
cd hass-crds

# Install CRDs
make install

# Deploy controller (builds and pushes image)
make deploy IMG=ghcr.io/your-org/hass-crds-controller:latest
```

### Environment Variables

The controller requires the following environment variables for MQTT connection:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MQTT_BROKER` | Yes | - | MQTT broker hostname |
| `MQTT_PORT` | No | `1883` | MQTT broker port |
| `MQTT_USERNAME` | No | - | MQTT username |
| `MQTT_PASSWORD` | No | - | MQTT password |
| `MQTT_CLIENT_ID` | No | auto-generated | MQTT client ID |
| `MQTT_USE_TLS` | No | `false` | Enable TLS (`true` or `1`) |

## Usage

### Basic Example: Button

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: garage-door-button
  namespace: home-automation
spec:
  name: "Garage Door"
  commandTopic: "home/garage/door/command"
  payloadPress: "TOGGLE"
  icon: "mdi:garage"
```

### Sensor with Device

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: living-room-temperature
  namespace: home-automation
spec:
  name: "Living Room Temperature"
  stateTopic: "home/living-room/temperature"
  unitOfMeasurement: "°C"
  deviceClass: "temperature"
  stateClass: "measurement"
  device:
    name: "Living Room Sensor Hub"
    identifiers:
      - "living-room-hub-001"
    manufacturer: "DIY"
    model: "ESP32 Sensor Hub"
```

### Switch with Availability

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSwitch
metadata:
  name: office-light
  namespace: home-automation
spec:
  name: "Office Light"
  commandTopic: "home/office/light/set"
  stateTopic: "home/office/light/state"
  payloadOn: "ON"
  payloadOff: "OFF"
  availability:
    - topic: "home/office/light/status"
      payloadAvailable: "online"
      payloadNotAvailable: "offline"
  availabilityMode: "all"
```

### Climate (Thermostat)

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTClimate
metadata:
  name: bedroom-thermostat
  namespace: home-automation
spec:
  name: "Bedroom Thermostat"
  modeCommandTopic: "home/bedroom/hvac/mode/set"
  modeStateTopic: "home/bedroom/hvac/mode/state"
  temperatureCommandTopic: "home/bedroom/hvac/temp/set"
  temperatureStateTopic: "home/bedroom/hvac/temp/state"
  currentTemperatureTopic: "home/bedroom/hvac/current"
  modes:
    - "off"
    - "heat"
    - "cool"
    - "auto"
  minTemp: 16
  maxTemp: 30
  tempStep: 0.5
```

### Light with Brightness and Color

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLight
metadata:
  name: kitchen-light
  namespace: home-automation
spec:
  name: "Kitchen Light"
  commandTopic: "home/kitchen/light/set"
  stateTopic: "home/kitchen/light/state"
  brightnessCommandTopic: "home/kitchen/light/brightness/set"
  brightnessStateTopic: "home/kitchen/light/brightness/state"
  brightnessScale: 255
  rgbCommandTopic: "home/kitchen/light/rgb/set"
  rgbStateTopic: "home/kitchen/light/rgb/state"
```

## How It Works

1. You create an MQTT entity CRD in Kubernetes
2. The controller detects the new resource
3. Controller builds an MQTT discovery payload
4. Payload is published to `homeassistant/<component>/<namespace>-<name>/config`
5. Home Assistant receives the discovery message and creates the entity
6. When you delete the CRD, the controller publishes an empty payload to remove the entity

### Topic Structure

Discovery topics follow the pattern:
```
homeassistant/<component>/<namespace>-<name>/config
```

For example, an `MQTTButton` named `garage-door` in namespace `home` publishes to:
```
homeassistant/button/home-garage-door/config
```

### Unique ID Generation

Each entity gets a unique ID automatically generated from `<namespace>-<name>`. You can override this with the `uniqueId` field in the spec.

## Development

### Build

```bash
make build          # Build controller binary
make docker-build   # Build Docker image
```

### Test

```bash
make test           # Unit tests
make test-e2e       # E2E tests (requires Docker, Kind)
make lint           # Run linter
```

### Generate

```bash
make generate       # Generate DeepCopy methods
make crds           # Generate CRD manifests (Python-based)
make manifests      # Generate RBAC manifests
```

## License

Apache License 2.0
