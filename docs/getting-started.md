# Getting Started

## Prerequisites

- **Kubernetes 1.24+** -- any conformant cluster (k3s, kind, GKE, EKS, etc.)
- **MQTT broker** -- accessible from the cluster (e.g. [Mosquitto](https://mosquitto.org/))
- **Home Assistant** -- with [MQTT integration](https://www.home-assistant.io/integrations/mqtt/) enabled and discovery turned on (enabled by default)
- **kubectl** -- configured to access your cluster

## Install CRDs

Apply the Custom Resource Definitions to your cluster:

```bash
kubectl apply -f https://raw.githubusercontent.com/spontus/hass-crds/main/config/crd/bases/
```

Verify the CRDs are installed:

```bash
kubectl get crds | grep mqtt.home-assistant.io
```

You should see output like:

```
mqttbuttons.mqtt.home-assistant.io        2025-01-15T10:00:00Z
mqttswitch.mqtt.home-assistant.io         2025-01-15T10:00:00Z
mqttlights.mqtt.home-assistant.io         2025-01-15T10:00:00Z
mqttsensors.mqtt.home-assistant.io        2025-01-15T10:00:00Z
mqttbinarysensors.mqtt.home-assistant.io  2025-01-15T10:00:00Z
mqttcovers.mqtt.home-assistant.io         2025-01-15T10:00:00Z
mqttclimates.mqtt.home-assistant.io       2025-01-15T10:00:00Z
mqttnumbers.mqtt.home-assistant.io        2025-01-15T10:00:00Z
mqttselects.mqtt.home-assistant.io        2025-01-15T10:00:00Z
mqtttexts.mqtt.home-assistant.io          2025-01-15T10:00:00Z
```

## Deploy the Controller

### Create a Namespace

```bash
kubectl create namespace hass-crds
```

### Configure MQTT Connection

Create a Secret with your MQTT broker credentials:

```bash
kubectl create secret generic mqtt-credentials \
  --namespace hass-crds \
  --from-literal=MQTT_HOST=mosquitto.default.svc.cluster.local \
  --from-literal=MQTT_PORT=1883 \
  --from-literal=MQTT_USERNAME=homeassistant \
  --from-literal=MQTT_PASSWORD=your-password-here
```

### Deploy

Apply the controller deployment:

```bash
kubectl apply -f https://raw.githubusercontent.com/spontus/hass-crds/main/config/deploy/controller.yaml
```

Or use the Helm chart:

```bash
helm install hass-crds hass-crds/hass-crds \
  --namespace hass-crds \
  --set mqtt.host=mosquitto.default.svc.cluster.local \
  --set mqtt.username=homeassistant \
  --set mqtt.password=your-password-here
```

### RBAC

The controller needs permission to watch and update CRD instances and their status subresources. The deployment manifests include the necessary RBAC resources. If you're deploying manually, ensure the controller's ServiceAccount has these permissions in its namespace:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: hass-crds-controller
  namespace: hass-crds
rules:
  - apiGroups: ["mqtt.home-assistant.io"]
    resources:
      - mqttbuttons
      - mqttswitch
      - mqttlights
      - mqttsensors
      - mqttbinarysensors
      - mqttcovers
      - mqttclimates
      - mqttnumbers
      - mqttselects
      - mqtttexts
    verbs: ["get", "list", "watch", "update", "patch"]
  - apiGroups: ["mqtt.home-assistant.io"]
    resources:
      - mqttbuttons/status
      - mqttswitch/status
      - mqttlights/status
      - mqttsensors/status
      - mqttbinarysensors/status
      - mqttcovers/status
      - mqttclimates/status
      - mqttnumbers/status
      - mqttselects/status
      - mqtttexts/status
    verbs: ["get", "update", "patch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]
```

## Verify the Installation

### Check Controller Logs

```bash
kubectl logs -n hass-crds deployment/hass-crds-controller
```

You should see:

```
INFO  Connected to MQTT broker at mosquitto.default.svc.cluster.local:1883
INFO  Watching mqttbuttons.mqtt.home-assistant.io/v1alpha1
INFO  Watching mqttswitch.mqtt.home-assistant.io/v1alpha1
...
INFO  Reconciliation loop started (interval: 60s)
```

### Create a Test Entity

Create a simple button to verify everything works:

```yaml
# test-button.yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: test-button
  namespace: hass-crds
spec:
  name: "Test Button"
  commandTopic: "hass-crds/test/button/command"
```

```bash
kubectl apply -f test-button.yaml
```

### Verify

1. Check the CRD status:

```bash
kubectl get mqttbutton test-button -n hass-crds -o yaml
```

Look for `.status.lastPublished` and `.status.discoveryTopic`.

2. Check Home Assistant -- the entity `button.test_button` should appear under the MQTT integration.

3. Clean up:

```bash
kubectl delete mqttbutton test-button -n hass-crds
```

The entity should disappear from Home Assistant.

## Next Steps

- Read the [CRD Reference](crds/README.md) for all supported entity types
- Browse [Examples](examples/README.md) for copy-pasteable manifests
- See [Controller](controller.md) for details on reconciliation, status, and error handling
