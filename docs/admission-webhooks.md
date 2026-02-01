# Admission Webhooks

The controller includes an optional validating admission webhook that catches configuration errors at apply time, before they reach the reconciliation loop.

## Enabling the Webhook

The webhook is included in the controller deployment but must be registered with a `ValidatingWebhookConfiguration`:

```bash
kubectl apply -f https://raw.githubusercontent.com/spontus/hass-crds/main/config/webhook/validating-webhook.yaml
```

The webhook requires TLS. The controller generates a self-signed certificate on startup and stores it in a Secret (`hass-crds-webhook-tls`). Alternatively, use [cert-manager](https://cert-manager.io/) for certificate management.

## Validation Rules

### Cross-Resource Uniqueness

| Check | Description |
|---|---|
| Duplicate `uniqueId` | Rejects a CR if another CR in the same namespace already uses the same `uniqueId` value |
| Duplicate discovery topic | Rejects a CR if it would publish to the same discovery topic as an existing CR |

### Field Validation

| Check | Description |
|---|---|
| Topic format | Topics must not be empty, must not contain null characters, and must not exceed 65535 bytes |
| Required topics | Entity types that require `commandTopic` or `stateTopic` are rejected if the field is missing and no `DEFAULT_TOPIC_TEMPLATE` is configured |
| Schema consistency | For `MQTTLight`, validates that fields match the selected `schema` (e.g. `brightness` is only valid with `schema: json`) |
| Enum values | Fields with a fixed set of values (e.g. `deviceClass`, `availabilityMode`, `schema`) are validated against allowed values |
| Numeric ranges | Fields like `qos` (0-2), `brightnessScale` (1+), `minTemp`/`maxTemp` (min < max) are range-checked |

### Device Reference Validation

| Check | Description |
|---|---|
| `MQTTDevice` exists | If `deviceRef` is used instead of an inline `device` block, the referenced `MQTTDevice` must exist in the same namespace |
| Device has identifiers | Warns if a device block has neither `identifiers` nor `connections` |

### Secret Reference Validation

| Check | Description |
|---|---|
| Secret exists | If a field uses a `secretRef`, the referenced Secret must exist in the same namespace |
| Key exists | The referenced key must exist within the Secret's data |

## Webhook Behavior

- **Failure policy**: `Fail` -- if the webhook is unreachable, CR creation/update is rejected. Set to `Ignore` if you prefer availability over validation.
- **Scope**: Namespace-scoped -- only validates CRs in namespaces with the label `mqtt.home-assistant.io/validate: "true"`
- **Side effects**: None -- the webhook only validates, it does not mutate resources

### Enabling Validation Per Namespace

Label the namespace to opt in:

```bash
kubectl label namespace hass-crds mqtt.home-assistant.io/validate=true
```

### Example Rejection

```bash
$ kubectl apply -f duplicate-sensor.yaml
Error from server (Forbidden): error when creating "duplicate-sensor.yaml":
  admission webhook "validate.mqtt.home-assistant.io" denied the request:
  uniqueId "living-room-temp" is already used by mqttsensor/temperature-sensor
  in namespace hass-crds
```

## Disabling the Webhook

Remove the webhook configuration to disable validation:

```bash
kubectl delete validatingwebhookconfiguration hass-crds-validating-webhook
```

CRs will still be validated by the CRD schema (structural validation), but cross-resource checks and advanced field validation will be skipped.
