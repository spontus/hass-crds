# Controller

The hass-crds controller is a Kubernetes controller that watches CRD instances and publishes MQTT discovery payloads to Home Assistant.

## Startup

On startup, the controller:

1. Reads configuration from environment variables
2. Connects to the MQTT broker
3. Registers watchers for all supported CRD types in its namespace
4. Starts the reconciliation loop

### Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `MQTT_HOST` | Yes | -- | MQTT broker hostname or IP |
| `MQTT_PORT` | No | `1883` | MQTT broker port |
| `MQTT_USERNAME` | No | -- | MQTT authentication username |
| `MQTT_PASSWORD` | No | -- | MQTT authentication password |
| `MQTT_DISCOVERY_PREFIX` | No | `homeassistant` | HA discovery topic prefix |
| `MQTT_TLS_ENABLED` | No | `false` | Enable TLS for MQTT connection |
| `MQTT_TLS_CA_CERT` | No | -- | Path to CA certificate for TLS |
| `MQTT_TLS_CLIENT_CERT` | No | -- | Path to client certificate for mutual TLS |
| `MQTT_TLS_CLIENT_KEY` | No | -- | Path to client key for mutual TLS |
| `RECONCILE_INTERVAL` | No | `60` | Seconds between periodic re-publishes |
| `LOG_LEVEL` | No | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `MQTT_TOPIC_PREFIX` | No | -- | Default prefix prepended to all entity topics (e.g. `devices/`). Entities can override with absolute topics. |
| `DEFAULT_TOPIC_TEMPLATE` | No | -- | Go template for auto-generating topics. Available variables: `{{.Namespace}}`, `{{.Name}}`, `{{.Component}}` (e.g. `{{.Component}}/{{.Namespace}}/{{.Name}}`) |
| `WATCH_NAMESPACE` | No | Controller's namespace | Namespace to watch for CRD instances. Defaults to the controller's own namespace. |

## Reconciliation Loop

The controller uses a standard Kubernetes reconciliation pattern:

### Event-Driven Reconciliation

When a CRD instance is created, updated, or deleted, the controller immediately reconciles:

1. **List** all instances of the CRD type in the current namespace
2. **Build** the MQTT discovery JSON payload from the CRD spec
   - Convert camelCase fields to snake_case keys
   - Auto-generate `unique_id` from `<namespace>-<name>` if not set
   - Derive the discovery topic: `<prefix>/<component>/<namespace>-<name>/config`
3. **Publish** the JSON payload to the discovery topic with `retain=true`
4. **Update** the CRD's status subresource

### Periodic Re-Publish

In addition to event-driven reconciliation, the controller periodically re-publishes all discovery payloads (default: every 60 seconds). This ensures:

- HA picks up entities after an MQTT broker restart
- Discovery payloads stay retained even if the broker loses them
- Any drift between the CRD spec and the published payload is corrected

The interval is configurable via the `RECONCILE_INTERVAL` environment variable.

### Per-Resource Rediscovery

Individual CRD instances can override the global re-publish interval using the `rediscoverInterval` field in their spec. See [Common Fields — Rediscovery](crds/common-fields.md#rediscovery) for details.

### Payload Construction

Given an MQTTButton:

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: restart-server
  namespace: default
spec:
  name: "Restart Server"
  commandTopic: "cmnd/server/restart"
  icon: "mdi:restart"
  device:
    name: "My Server"
    identifiers:
      - "server-01"
```

The controller publishes the following JSON to `homeassistant/button/default-restart-server/config`:

```json
{
  "name": "Restart Server",
  "command_topic": "cmnd/server/restart",
  "icon": "mdi:restart",
  "unique_id": "default-restart-server",
  "device": {
    "name": "My Server",
    "identifiers": ["server-01"]
  }
}
```

## Topic Defaults

When `MQTT_TOPIC_PREFIX` or `DEFAULT_TOPIC_TEMPLATE` is set on the controller, topic fields in CRD specs can use shorter relative paths.

### Topic Prefix

If `MQTT_TOPIC_PREFIX` is set to `devices/`, a CRD with `commandTopic: "living-room/lamp/set"` publishes to `devices/living-room/lamp/set`.

Topics that begin with `/` are treated as absolute and bypass the prefix:

```yaml
spec:
  # This uses the prefix: resolves to "devices/lamp/set"
  commandTopic: "lamp/set"
  # This bypasses the prefix: resolves to "cmnd/lamp/POWER"
  stateTopic: "/cmnd/lamp/POWER"
```

### Topic Template

If `DEFAULT_TOPIC_TEMPLATE` is set, any CRD that omits its required topic fields will have them auto-generated from the template. For example, with `DEFAULT_TOPIC_TEMPLATE={{.Component}}/{{.Namespace}}/{{.Name}}`:

- An `MQTTSwitch` named `desk-lamp` in namespace `office` gets `commandTopic` set to `switch/office/desk-lamp/set` and `stateTopic` set to `switch/office/desk-lamp/state`

Explicitly set topics always take priority over the template.

## Deletion

When a CRD instance is deleted, the controller's behavior depends on the finalizer policy.

### Default: Clean Removal

By default, the controller removes entities from Home Assistant on CR deletion:

1. Kubernetes marks the object for deletion but the **finalizer** (`mqtt.home-assistant.io/discovery`) prevents actual removal
2. The controller detects the deletion timestamp and:
   - Publishes an **empty payload** (`""`) to the discovery topic with `retain=true`
   - This tells Home Assistant to remove the entity
3. The controller removes the finalizer from the object
4. Kubernetes completes the deletion

### Orphan Mode

In some cases you may want to delete the CR without removing the entity from Home Assistant (e.g. during controller migration or maintenance). Annotate the resource before deleting:

```bash
kubectl annotate mqttbutton my-button mqtt.home-assistant.io/deletion-policy=orphan
kubectl delete mqttbutton my-button
```

When `mqtt.home-assistant.io/deletion-policy: orphan` is set, the controller skips publishing the empty payload and simply removes the finalizer. The entity remains in Home Assistant with its last-known configuration.

| Annotation Value | Behavior |
|---|---|
| `cleanup` (default) | Publish empty payload, remove entity from HA |
| `orphan` | Skip cleanup, leave entity in HA |

## Status Subresource

Each CRD instance has a `.status` subresource updated by the controller:

| Field | Type | Description |
|---|---|---|
| `.status.lastPublished` | `string` (RFC 3339 timestamp) | When the discovery payload was last published |
| `.status.discoveryTopic` | `string` | The MQTT topic the payload was published to |
| `.status.conditions` | `[]Condition` | Standard Kubernetes conditions |

### Conditions

| Type | Description |
|---|---|
| `Published` | `True` when the discovery payload has been successfully published |
| `MQTTConnected` | `True` when the controller has an active MQTT connection |

Example status:

```yaml
status:
  lastPublished: "2025-01-15T10:30:00Z"
  discoveryTopic: "homeassistant/button/default-restart-server/config"
  conditions:
    - type: Published
      status: "True"
      lastTransitionTime: "2025-01-15T10:30:00Z"
      message: "Discovery payload published successfully"
    - type: MQTTConnected
      status: "True"
      lastTransitionTime: "2025-01-15T10:00:00Z"
```

## Error Handling

### MQTT Connection Failures

- The controller uses **exponential backoff** for reconnection (starting at 1s, capped at 60s)
- While disconnected, the controller continues to reconcile CRDs but marks `MQTTConnected=False` in status conditions
- On reconnection, all discovery payloads are immediately re-published

### Kubernetes API Errors

- Transient errors (network issues, API server overload) trigger a requeue with backoff
- The controller emits Kubernetes **Events** on the CRD instance for notable operations:
  - `Normal/Published` -- discovery payload published successfully
  - `Warning/PublishFailed` -- failed to publish (includes error detail)
  - `Normal/Deleted` -- empty payload published for cleanup

### Invalid CRD Specs

- Schema validation catches most errors at admission time
- Runtime errors (e.g. topics that are too long) are reported via conditions and events

## Dry Run

Annotate a CRD instance with `mqtt.home-assistant.io/dry-run: "true"` to make the controller log the discovery JSON it would publish without actually publishing to MQTT. This is useful for debugging payloads before they reach Home Assistant.

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: test-sensor
  annotations:
    mqtt.home-assistant.io/dry-run: "true"
spec:
  name: "Test Sensor"
  stateTopic: "sensors/test/state"
```

When dry-run is enabled:

- The controller logs the full JSON payload at `info` level
- The `Published` condition is set to `False` with reason `DryRun`
- No MQTT messages are published
- The `.status.discoveryTopic` is still populated so you can see where it *would* publish

Remove the annotation and the controller publishes on the next reconciliation.

## Multi-Tenancy

The controller is **namespace-scoped** by design. Each namespace with its own controller deployment acts as an independent tenant.

### Deployment Model

```
Namespace: team-a               Namespace: team-b
┌──────────────────────┐        ┌──────────────────────┐
│ controller (team-a)  │        │ controller (team-b)  │
│  MQTT_HOST=broker-a  │        │  MQTT_HOST=broker-b  │
│  MQTT_TOPIC_PREFIX=  │        │  MQTT_TOPIC_PREFIX=  │
│    team-a/           │        │    team-b/           │
│  MQTT_DISCOVERY_     │        │  MQTT_DISCOVERY_     │
│    PREFIX=ha         │        │    PREFIX=ha         │
├──────────────────────┤        ├──────────────────────┤
│ MQTTSensor/temp      │        │ MQTTLight/desk-lamp  │
│ MQTTSwitch/relay     │        │ MQTTFan/ceiling-fan  │
└──────────────────────┘        └──────────────────────┘
         │                               │
         v                               v
    MQTT Broker A                   MQTT Broker B
         │                               │
         v                               v
  Home Assistant A               Home Assistant B
```

Each controller has its own:

- **MQTT broker connection** -- different teams can use different brokers or share one
- **Discovery prefix** -- isolates HA instances from each other
- **Topic prefix** -- prevents topic collisions between tenants
- **RBAC** -- controller only needs permissions in its own namespace

### Shared Broker

Multiple controllers can share a single MQTT broker. Use distinct `MQTT_DISCOVERY_PREFIX` values to target different HA instances, or distinct `MQTT_TOPIC_PREFIX` values to keep topics separated:

```yaml
# team-a controller
env:
  - name: MQTT_HOST
    value: "shared-broker.mqtt.svc"
  - name: MQTT_DISCOVERY_PREFIX
    value: "homeassistant"
  - name: MQTT_TOPIC_PREFIX
    value: "team-a/"

# team-b controller
env:
  - name: MQTT_HOST
    value: "shared-broker.mqtt.svc"
  - name: MQTT_DISCOVERY_PREFIX
    value: "homeassistant"
  - name: MQTT_TOPIC_PREFIX
    value: "team-b/"
```
