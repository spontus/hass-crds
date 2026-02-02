# Security Considerations

Best practices for securing your hass-crds deployment.

## Secrets Management

### Never Store Secrets in CRD Specs

Sensitive values like alarm codes, lock codes, and passwords should never be stored in plaintext in CRD specs.

**Bad**:
```yaml
spec:
  code: "1234"  # Exposed in etcd, logs, and kubectl output
```

**Good**:
```yaml
spec:
  code:
    secretRef:
      name: alarm-codes
      key: panel-code
```

See [Common Fields - Secret References](crds/common-fields.md#secret-references) for details.

### MQTT Credentials

Store MQTT credentials in Kubernetes Secrets:

```bash
kubectl create secret generic mqtt-credentials \
  --namespace hass-crds \
  --from-literal=MQTT_USERNAME=homeassistant \
  --from-literal=MQTT_PASSWORD='<strong-password>'
```

Reference in the controller deployment:
```yaml
env:
  - name: MQTT_USERNAME
    valueFrom:
      secretKeyRef:
        name: mqtt-credentials
        key: MQTT_USERNAME
  - name: MQTT_PASSWORD
    valueFrom:
      secretKeyRef:
        name: mqtt-credentials
        key: MQTT_PASSWORD
```

### Secret Rotation

When rotating secrets:
1. Update the Kubernetes Secret
2. The controller automatically picks up changes on next reconciliation
3. CRDs referencing the secret are re-published with new values

## Network Security

### MQTT TLS

Always use TLS for MQTT connections in production:

```yaml
env:
  - name: MQTT_TLS_ENABLED
    value: "true"
  - name: MQTT_TLS_CA_CERT
    value: "/etc/mqtt/certs/ca.crt"
```

See [Controller - TLS Configuration](controller.md#tls-configuration) for setup details.

### Network Policies

Restrict controller network access to only the MQTT broker:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: hass-crds-controller
  namespace: hass-crds
spec:
  podSelector:
    matchLabels:
      app: hass-crds-controller
  policyTypes:
    - Egress
  egress:
    # Allow MQTT broker access
    - to:
        - namespaceSelector:
            matchLabels:
              name: mqtt
          podSelector:
            matchLabels:
              app: mosquitto
      ports:
        - protocol: TCP
          port: 8883
    # Allow Kubernetes API access
    - to:
        - ipBlock:
            cidr: 10.0.0.1/32  # Kubernetes API server
      ports:
        - protocol: TCP
          port: 443
    # Allow DNS
    - to: []
      ports:
        - protocol: UDP
          port: 53
```

### MQTT Broker ACLs

Configure your MQTT broker to restrict topic access:

**Mosquitto example** (`/etc/mosquitto/acl`):
```
# hass-crds controller can publish to discovery topics
user hass-crds
topic write homeassistant/#

# hass-crds controller can publish/subscribe to entity topics
user hass-crds
topic readwrite devices/#
```

## RBAC

### Principle of Least Privilege

The controller only needs permissions in its own namespace:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role  # Not ClusterRole
metadata:
  name: hass-crds-controller
  namespace: hass-crds
rules:
  - apiGroups: ["mqtt.home-assistant.io"]
    resources: ["*"]
    verbs: ["get", "list", "watch", "update", "patch"]
  - apiGroups: ["mqtt.home-assistant.io"]
    resources: ["*/status"]
    verbs: ["get", "update", "patch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get"]  # Only get, not list or watch
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]
```

### Multi-Tenant Isolation

Each team's namespace should have its own controller with isolated RBAC:

```yaml
# team-a namespace
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: hass-crds-controller
  namespace: team-a
subjects:
  - kind: ServiceAccount
    name: hass-crds-controller
    namespace: team-a  # SA in same namespace
roleRef:
  kind: Role
  name: hass-crds-controller
  apiGroup: rbac.authorization.k8s.io
```

## Audit Logging

### Controller Logs

The controller logs all publish operations. Enable debug logging for detailed output:

```yaml
env:
  - name: LOG_LEVEL
    value: "debug"
```

Sensitive values (from secretRef) are never logged.

### Kubernetes Audit

Enable Kubernetes audit logging for CRD operations:

```yaml
# audit-policy.yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
  - level: RequestResponse
    resources:
      - group: "mqtt.home-assistant.io"
        resources: ["*"]
```

## Admission Webhooks

Enable the validating webhook for additional security checks:

- Prevents duplicate `uniqueId` values (which could cause entity hijacking)
- Validates Secret references exist before accepting CRs
- Enforces schema consistency

See [Admission Webhooks](admission-webhooks.md) for setup.

## Supply Chain Security

### Image Verification

Verify controller image signatures (when available):

```bash
cosign verify ghcr.io/spontus/hass-crds-controller:latest
```

### Image Pull Policy

Use specific image tags and `imagePullPolicy: IfNotPresent` in production:

```yaml
containers:
  - name: controller
    image: ghcr.io/spontus/hass-crds-controller:v1.0.0  # Pin version
    imagePullPolicy: IfNotPresent
```

## Reporting Security Issues

To report a security vulnerability, please email security@example.com or open a private security advisory on GitHub.

Do not open public issues for security vulnerabilities.
