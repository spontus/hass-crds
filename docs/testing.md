# Testing Strategy

This document describes the testing approach for hass-crds, covering unit tests, integration tests, and end-to-end tests.

## Test Pyramid

```
        ┌───────────┐
        │   E2E     │  Few, slow, high confidence
        │  Tests    │
        ├───────────┤
        │Integration│  Some, moderate speed
        │  Tests    │
        ├───────────┤
        │   Unit    │  Many, fast, isolated
        │   Tests   │
        └───────────┘
```

| Level | What It Tests | Speed | Dependencies |
|-------|---------------|-------|--------------|
| Unit | Individual functions, payload construction | Fast (ms) | None |
| Integration | Controller + MQTT broker interaction | Medium (s) | MQTT broker |
| E2E | Full flow: CRD → Controller → MQTT → Home Assistant | Slow (min) | Full stack |

## Unit Tests

Unit tests validate individual components in isolation without external dependencies.

### What to Unit Test

| Component | Test Cases |
|-----------|------------|
| **Payload Builder** | camelCase → snake_case conversion, field mapping, JSON structure |
| **Topic Derivation** | Discovery topic generation, prefix handling, absolute vs relative topics |
| **Unique ID Generation** | Auto-generation from namespace/name, custom uniqueId passthrough |
| **Field Validation** | Required fields, enum values, numeric ranges |
| **Template Rendering** | Topic templates with namespace/name/component variables |

### Running Unit Tests

```bash
go test ./... -v -short
```

### Example: Payload Builder Test

```go
// internal/payload/builder_test.go
package payload

import (
    "testing"
    "github.com/stretchr/testify/assert"
    v1alpha1 "github.com/spontus/hass-crds/api/v1alpha1"
)

func TestBuildButtonPayload(t *testing.T) {
    button := &v1alpha1.MQTTButton{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "restart-server",
            Namespace: "default",
        },
        Spec: v1alpha1.MQTTButtonSpec{
            Name:         "Restart Server",
            CommandTopic: "cmnd/server/restart",
            PayloadPress: "RESTART",
            Icon:         "mdi:restart",
        },
    }

    payload, err := BuildButtonPayload(button)

    assert.NoError(t, err)
    assert.Equal(t, "Restart Server", payload["name"])
    assert.Equal(t, "cmnd/server/restart", payload["command_topic"])
    assert.Equal(t, "RESTART", payload["payload_press"])
    assert.Equal(t, "mdi:restart", payload["icon"])
    assert.Equal(t, "default-restart-server", payload["unique_id"])
}

func TestCamelToSnakeCase(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"commandTopic", "command_topic"},
        {"stateTopic", "state_topic"},
        {"uniqueId", "unique_id"},
        {"jsonAttributesTopic", "json_attributes_topic"},
        {"rgbCommandTopic", "rgb_command_topic"},
    }

    for _, tc := range tests {
        t.Run(tc.input, func(t *testing.T) {
            result := camelToSnake(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### Example: Topic Derivation Test

```go
// internal/topic/topic_test.go
package topic

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDiscoveryTopic(t *testing.T) {
    tests := []struct {
        prefix    string
        component string
        namespace string
        name      string
        expected  string
    }{
        {"homeassistant", "button", "default", "restart", "homeassistant/button/default-restart/config"},
        {"ha", "sensor", "monitoring", "cpu-temp", "ha/sensor/monitoring-cpu-temp/config"},
        {"homeassistant", "binary_sensor", "security", "motion", "homeassistant/binary_sensor/security-motion/config"},
    }

    for _, tc := range tests {
        t.Run(tc.expected, func(t *testing.T) {
            result := DiscoveryTopic(tc.prefix, tc.component, tc.namespace, tc.name)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func TestApplyTopicPrefix(t *testing.T) {
    tests := []struct {
        prefix   string
        topic    string
        expected string
    }{
        {"devices/", "lamp/set", "devices/lamp/set"},
        {"devices/", "/absolute/topic", "/absolute/topic"},  // Absolute bypasses prefix
        {"", "lamp/set", "lamp/set"},                        // No prefix
    }

    for _, tc := range tests {
        t.Run(tc.topic, func(t *testing.T) {
            result := ApplyPrefix(tc.prefix, tc.topic)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### Example: Validation Test

```go
// internal/validation/validation_test.go
package validation

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidateQoS(t *testing.T) {
    assert.NoError(t, ValidateQoS(0))
    assert.NoError(t, ValidateQoS(1))
    assert.NoError(t, ValidateQoS(2))
    assert.Error(t, ValidateQoS(-1))
    assert.Error(t, ValidateQoS(3))
}

func TestValidateDeviceClass(t *testing.T) {
    // Button device classes
    assert.NoError(t, ValidateButtonDeviceClass("identify"))
    assert.NoError(t, ValidateButtonDeviceClass("restart"))
    assert.NoError(t, ValidateButtonDeviceClass("update"))
    assert.NoError(t, ValidateButtonDeviceClass(""))  // Optional
    assert.Error(t, ValidateButtonDeviceClass("invalid"))
}
```

## Integration Tests

Integration tests verify the controller works correctly with a real MQTT broker.

### Prerequisites

- Docker (for running MQTT broker)
- kubectl configured for a test cluster (kind, k3d, or minikube)

### Test Environment Setup

```bash
# Start MQTT broker
docker run -d --name mqtt-test -p 1883:1883 eclipse-mosquitto:2 \
  mosquitto -c /mosquitto-no-auth.conf

# Install CRDs
kubectl apply -f config/crd/crds.yaml

# Create test namespace
kubectl create namespace hass-crds-test
```

### Running Integration Tests

```bash
export MQTT_HOST=localhost
export MQTT_PORT=1883
export TEST_NAMESPACE=hass-crds-test

go test ./... -v -tags=integration
```

### Example: Reconciliation Integration Test

```go
//go:build integration

// internal/controller/button_controller_integration_test.go
package controller

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    v1alpha1 "github.com/spontus/hass-crds/api/v1alpha1"
)

func TestButtonReconciliation(t *testing.T) {
    ctx := context.Background()

    // Setup MQTT client to verify published messages
    mqttClient := newTestMQTTClient(t)
    defer mqttClient.Disconnect(250)

    // Subscribe to discovery topic
    received := make(chan []byte, 1)
    token := mqttClient.Subscribe("homeassistant/button/+/config", 0, func(_ mqtt.Client, msg mqtt.Message) {
        received <- msg.Payload()
    })
    require.NoError(t, token.Error())

    // Create MQTTButton resource
    button := &v1alpha1.MQTTButton{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-button",
            Namespace: testNamespace,
        },
        Spec: v1alpha1.MQTTButtonSpec{
            Name:         "Test Button",
            CommandTopic: "test/button/command",
        },
    }

    err := k8sClient.Create(ctx, button)
    require.NoError(t, err)
    defer k8sClient.Delete(ctx, button)

    // Wait for discovery message
    select {
    case payload := <-received:
        var discovery map[string]interface{}
        err := json.Unmarshal(payload, &discovery)
        require.NoError(t, err)

        assert.Equal(t, "Test Button", discovery["name"])
        assert.Equal(t, "test/button/command", discovery["command_topic"])
        assert.Contains(t, discovery["unique_id"], "test-button")

    case <-time.After(10 * time.Second):
        t.Fatal("Timeout waiting for discovery message")
    }

    // Verify status was updated
    var updated v1alpha1.MQTTButton
    err = k8sClient.Get(ctx, client.ObjectKeyFromObject(button), &updated)
    require.NoError(t, err)

    assert.NotEmpty(t, updated.Status.LastPublished)
    assert.Contains(t, updated.Status.DiscoveryTopic, "homeassistant/button")
}
```

### Example: Deletion Integration Test

```go
//go:build integration

func TestButtonDeletion(t *testing.T) {
    ctx := context.Background()

    mqttClient := newTestMQTTClient(t)
    defer mqttClient.Disconnect(250)

    // Track messages on discovery topic
    messages := make(chan []byte, 10)
    expectedTopic := fmt.Sprintf("homeassistant/button/%s-delete-test/config", testNamespace)

    mqttClient.Subscribe(expectedTopic, 0, func(_ mqtt.Client, msg mqtt.Message) {
        messages <- msg.Payload()
    })

    // Create button
    button := &v1alpha1.MQTTButton{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "delete-test",
            Namespace: testNamespace,
        },
        Spec: v1alpha1.MQTTButtonSpec{
            Name:         "Delete Test",
            CommandTopic: "test/delete/command",
        },
    }

    err := k8sClient.Create(ctx, button)
    require.NoError(t, err)

    // Wait for initial publish
    select {
    case payload := <-messages:
        assert.NotEmpty(t, payload)
    case <-time.After(10 * time.Second):
        t.Fatal("Timeout waiting for initial publish")
    }

    // Delete the button
    err = k8sClient.Delete(ctx, button)
    require.NoError(t, err)

    // Wait for empty payload (cleanup)
    select {
    case payload := <-messages:
        assert.Empty(t, payload, "Expected empty payload for cleanup")
    case <-time.After(10 * time.Second):
        t.Fatal("Timeout waiting for cleanup message")
    }
}
```

### Example: Secret Reference Integration Test

```go
//go:build integration

func TestSecretRefResolution(t *testing.T) {
    ctx := context.Background()

    // Create secret
    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-codes",
            Namespace: testNamespace,
        },
        StringData: map[string]string{
            "alarm-code": "1234",
        },
    }
    err := k8sClient.Create(ctx, secret)
    require.NoError(t, err)
    defer k8sClient.Delete(ctx, secret)

    // Create alarm panel with secretRef
    alarm := &v1alpha1.MQTTAlarmControlPanel{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-alarm",
            Namespace: testNamespace,
        },
        Spec: v1alpha1.MQTTAlarmControlPanelSpec{
            Name:         "Test Alarm",
            CommandTopic: "alarm/set",
            StateTopic:   "alarm/state",
            Code: v1alpha1.StringOrSecretRef{
                SecretRef: &v1alpha1.SecretKeyRef{
                    Name: "test-codes",
                    Key:  "alarm-code",
                },
            },
        },
    }

    // ... verify the resolved value appears in MQTT payload
}
```

## End-to-End Tests

E2E tests verify the complete flow from CRD creation to entity appearing in Home Assistant.

### Prerequisites

- Running Home Assistant instance with MQTT integration
- Long-lived access token for HA API
- MQTT broker connected to both controller and HA

### Environment Setup

```bash
export HA_URL=http://homeassistant.local:8123
export HA_TOKEN=<long-lived-access-token>
export MQTT_HOST=mqtt.local
export MQTT_PORT=1883
```

### Running E2E Tests

```bash
go test ./... -v -tags=e2e -timeout=5m
```

### Example: E2E Entity Creation Test

```go
//go:build e2e

// test/e2e/entity_creation_test.go
package e2e

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestButtonAppearsInHomeAssistant(t *testing.T) {
    ctx := context.Background()

    // Create MQTTButton via kubectl
    button := &v1alpha1.MQTTButton{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "e2e-test-button",
            Namespace: "hass-crds",
        },
        Spec: v1alpha1.MQTTButtonSpec{
            Name:         "E2E Test Button",
            CommandTopic: "e2e/button/command",
            DeviceClass:  "restart",
        },
    }

    err := k8sClient.Create(ctx, button)
    require.NoError(t, err)
    defer k8sClient.Delete(ctx, button)

    // Wait for entity to appear in Home Assistant
    entityID := "button.e2e_test_button"
    var entity HAEntity

    require.Eventually(t, func() bool {
        entity, err = getHAEntity(entityID)
        return err == nil
    }, 30*time.Second, 1*time.Second, "Entity did not appear in HA")

    assert.Equal(t, "E2E Test Button", entity.Attributes.FriendlyName)
    assert.Equal(t, "restart", entity.Attributes.DeviceClass)
}

func TestButtonDisappearsOnDeletion(t *testing.T) {
    ctx := context.Background()

    // Create button
    button := &v1alpha1.MQTTButton{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "e2e-delete-test",
            Namespace: "hass-crds",
        },
        Spec: v1alpha1.MQTTButtonSpec{
            Name:         "Delete Test",
            CommandTopic: "e2e/delete/command",
        },
    }

    err := k8sClient.Create(ctx, button)
    require.NoError(t, err)

    // Wait for entity to appear
    entityID := "button.delete_test"
    require.Eventually(t, func() bool {
        _, err := getHAEntity(entityID)
        return err == nil
    }, 30*time.Second, 1*time.Second)

    // Delete the button
    err = k8sClient.Delete(ctx, button)
    require.NoError(t, err)

    // Wait for entity to disappear
    require.Eventually(t, func() bool {
        _, err := getHAEntity(entityID)
        return err != nil  // Entity should not be found
    }, 30*time.Second, 1*time.Second, "Entity still exists after deletion")
}

// Helper: Get entity from Home Assistant API
func getHAEntity(entityID string) (HAEntity, error) {
    url := fmt.Sprintf("%s/api/states/%s", os.Getenv("HA_URL"), entityID)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+os.Getenv("HA_TOKEN"))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return HAEntity{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
        return HAEntity{}, fmt.Errorf("entity not found")
    }

    var entity HAEntity
    json.NewDecoder(resp.Body).Decode(&entity)
    return entity, nil
}

type HAEntity struct {
    EntityID   string `json:"entity_id"`
    State      string `json:"state"`
    Attributes struct {
        FriendlyName string `json:"friendly_name"`
        DeviceClass  string `json:"device_class"`
    } `json:"attributes"`
}
```

## Test Coverage

### Generating Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Targets

| Package | Target Coverage |
|---------|-----------------|
| `internal/payload` | 90%+ |
| `internal/topic` | 90%+ |
| `internal/validation` | 85%+ |
| `internal/controller` | 80%+ |

## CI/CD Integration

### GitHub Actions Example

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test ./... -v -short -coverprofile=coverage.out
      - name: Upload coverage
        uses: codecov/codecov-action@v4

  integration-tests:
    runs-on: ubuntu-latest
    services:
      mqtt:
        image: eclipse-mosquitto:2
        ports:
          - 1883:1883
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Create kind cluster
        uses: helm/kind-action@v1
      - name: Install CRDs
        run: kubectl apply -f config/crd/crds.yaml
      - name: Run integration tests
        env:
          MQTT_HOST: localhost
          MQTT_PORT: 1883
        run: go test ./... -v -tags=integration
```

## Test Utilities

### Test Fixtures

Store reusable test fixtures in `testdata/`:

```
testdata/
├── crds/
│   ├── valid-button.yaml
│   ├── valid-sensor.yaml
│   └── invalid-missing-topic.yaml
├── payloads/
│   ├── expected-button.json
│   └── expected-sensor.json
└── secrets/
    └── test-secret.yaml
```

### Test Helpers

```go
// internal/testutil/helpers.go
package testutil

import (
    "os"
    "path/filepath"
    "testing"

    "sigs.k8s.io/yaml"
)

func LoadFixture[T any](t *testing.T, path string) T {
    t.Helper()

    data, err := os.ReadFile(filepath.Join("testdata", path))
    if err != nil {
        t.Fatalf("Failed to load fixture %s: %v", path, err)
    }

    var obj T
    if err := yaml.Unmarshal(data, &obj); err != nil {
        t.Fatalf("Failed to unmarshal fixture %s: %v", path, err)
    }

    return obj
}

func NewTestMQTTClient(t *testing.T) mqtt.Client {
    t.Helper()

    opts := mqtt.NewClientOptions().
        AddBroker(fmt.Sprintf("tcp://%s:%s",
            os.Getenv("MQTT_HOST"),
            os.Getenv("MQTT_PORT"))).
        SetClientID(fmt.Sprintf("test-%d", time.Now().UnixNano()))

    client := mqtt.NewClient(opts)
    token := client.Connect()
    token.Wait()

    if err := token.Error(); err != nil {
        t.Fatalf("Failed to connect to MQTT: %v", err)
    }

    return client
}
```

## Mocking

### Mock MQTT Client

```go
// internal/mqtt/mock_client.go
package mqtt

type MockClient struct {
    Published []PublishedMessage
    Connected bool
}

type PublishedMessage struct {
    Topic   string
    Payload []byte
    QoS     byte
    Retain  bool
}

func (m *MockClient) Publish(topic string, qos byte, retain bool, payload interface{}) error {
    m.Published = append(m.Published, PublishedMessage{
        Topic:   topic,
        Payload: payload.([]byte),
        QoS:     qos,
        Retain:  retain,
    })
    return nil
}

func (m *MockClient) IsConnected() bool {
    return m.Connected
}
```

### Using Mocks in Tests

```go
func TestReconcilerPublishesCorrectPayload(t *testing.T) {
    mockMQTT := &mqtt.MockClient{Connected: true}

    reconciler := &ButtonReconciler{
        Client:     k8sClient,
        MQTTClient: mockMQTT,
    }

    // Trigger reconciliation...

    require.Len(t, mockMQTT.Published, 1)
    assert.Equal(t, "homeassistant/button/default-test/config", mockMQTT.Published[0].Topic)
    assert.True(t, mockMQTT.Published[0].Retain)
}
```

---

## See Also

- [Controller](controller.md) - Controller behavior and configuration
- [Contributing](../CONTRIBUTING.md) - Development setup and guidelines
- [Architecture](architecture.md) - Design decisions
