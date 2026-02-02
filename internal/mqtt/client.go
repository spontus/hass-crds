/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mqtt

import (
	"context"
	"fmt"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-logr/logr"
)

const (
	// DefaultKeepAlive is the default MQTT keep-alive interval.
	DefaultKeepAlive = 30 * time.Second

	// DefaultConnectTimeout is the default timeout for initial connection.
	DefaultConnectTimeout = 30 * time.Second

	// DefaultWriteTimeout is the default timeout for publish operations.
	DefaultWriteTimeout = 10 * time.Second

	// DefaultMaxReconnectInterval is the maximum interval between reconnect attempts.
	DefaultMaxReconnectInterval = 5 * time.Minute

	// DefaultReconnectWaitTimeout is how long Publish waits for reconnection.
	DefaultReconnectWaitTimeout = 30 * time.Second
)

// Client defines the interface for MQTT operations.
type Client interface {
	Connect(ctx context.Context) error
	Disconnect()
	Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error
	IsConnected() bool
	WaitForConnection(ctx context.Context) error
}

// PahoClient wraps the Paho MQTT client.
type PahoClient struct {
	client        pahomqtt.Client
	config        *Config
	log           logr.Logger
	mu            sync.RWMutex
	disconnecting bool
}

// NewClient creates a new MQTT client with the given configuration.
func NewClient(config *Config, log logr.Logger) *PahoClient {
	return &PahoClient{
		config: config,
		log:    log.WithName("mqtt-client"),
	}
}

// Connect establishes the MQTT connection.
func (c *PahoClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.disconnecting = false

	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(c.config.BrokerURL())
	opts.SetClientID(c.config.ClientID)
	opts.SetCleanSession(true)
	opts.SetKeepAlive(DefaultKeepAlive)
	opts.SetWriteTimeout(DefaultWriteTimeout)
	opts.SetConnectTimeout(DefaultConnectTimeout)

	// Auto-reconnect settings
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(DefaultMaxReconnectInterval)

	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
		opts.SetPassword(c.config.Password)
	}

	opts.SetConnectionLostHandler(func(client pahomqtt.Client, err error) {
		c.log.Error(err, "MQTT connection lost, will auto-reconnect", "broker", c.config.BrokerURL())
	})

	opts.SetOnConnectHandler(func(client pahomqtt.Client) {
		c.log.Info("MQTT connected", "broker", c.config.BrokerURL())
	})

	opts.SetReconnectingHandler(func(client pahomqtt.Client, opts *pahomqtt.ClientOptions) {
		c.log.Info("MQTT attempting reconnection", "broker", c.config.BrokerURL())
	})

	c.client = pahomqtt.NewClient(opts)

	token := c.client.Connect()

	// Wait for connection with context timeout
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
		}
	}

	c.log.Info("MQTT client connected", "broker", c.config.BrokerURL())
	return nil
}

// Disconnect closes the MQTT connection.
func (c *PahoClient) Disconnect() {
	c.mu.Lock()
	c.disconnecting = true

	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(1000) // 1 second timeout
		c.log.Info("MQTT client disconnected")
	}
	c.mu.Unlock()
}

// Publish sends a message to the specified topic.
// If not connected, it waits for reconnection up to DefaultReconnectWaitTimeout.
func (c *PahoClient) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
	// Wait for connection if needed
	if err := c.WaitForConnection(ctx); err != nil {
		return err
	}

	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	token := client.Publish(topic, qos, retain, payload)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return fmt.Errorf("failed to publish to %s: %w", topic, token.Error())
		}
	}

	c.log.V(1).Info("Published MQTT message", "topic", topic, "retain", retain, "qos", qos)
	return nil
}

// WaitForConnection waits for the MQTT client to be connected.
// Returns immediately if already connected, otherwise waits for reconnection.
func (c *PahoClient) WaitForConnection(ctx context.Context) error {
	// Quick check without lock
	c.mu.RLock()
	client := c.client
	disconnecting := c.disconnecting
	c.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("MQTT client not initialized")
	}

	if disconnecting {
		return fmt.Errorf("MQTT client is disconnecting")
	}

	// Already connected - fast path
	if client.IsConnected() {
		return nil
	}

	c.log.Info("Waiting for MQTT reconnection before publish")

	// Poll for reconnection with timeout
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.NewTimer(DefaultReconnectWaitTimeout)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timeout.C:
			return fmt.Errorf("timeout waiting for MQTT reconnection after %v", DefaultReconnectWaitTimeout)

		case <-ticker.C:
			c.mu.RLock()
			connected := c.client != nil && c.client.IsConnected()
			disconnecting := c.disconnecting
			c.mu.RUnlock()

			if disconnecting {
				return fmt.Errorf("MQTT client is disconnecting")
			}

			if connected {
				c.log.Info("MQTT reconnected, proceeding with publish")
				return nil
			}
		}
	}
}

// IsConnected returns whether the client is connected.
func (c *PahoClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.client != nil && c.client.IsConnected()
}

// MockClient is a mock MQTT client for testing.
type MockClient struct {
	connected      bool
	publishedMsgs  []PublishedMessage
	mu             sync.Mutex
	publishErr     error
	connectErr     error
}

// PublishedMessage records a published message for testing.
type PublishedMessage struct {
	Topic   string
	Payload []byte
	QoS     byte
	Retain  bool
}

// NewMockClient creates a new mock MQTT client.
func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connectErr != nil {
		return m.connectErr
	}
	m.connected = true
	return nil
}

func (m *MockClient) Disconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
}

func (m *MockClient) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.publishErr != nil {
		return m.publishErr
	}

	m.publishedMsgs = append(m.publishedMsgs, PublishedMessage{
		Topic:   topic,
		Payload: payload,
		QoS:     qos,
		Retain:  retain,
	})
	return nil
}

func (m *MockClient) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.connected
}

func (m *MockClient) WaitForConnection(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connectErr != nil {
		return m.connectErr
	}
	if !m.connected {
		return fmt.Errorf("mock client not connected")
	}
	return nil
}

// GetPublishedMessages returns all published messages for testing.
func (m *MockClient) GetPublishedMessages() []PublishedMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]PublishedMessage, len(m.publishedMsgs))
	copy(result, m.publishedMsgs)
	return result
}

// SetPublishError sets an error to return on publish.
func (m *MockClient) SetPublishError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishErr = err
}

// SetConnectError sets an error to return on connect.
func (m *MockClient) SetConnectError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectErr = err
}

// ClearMessages clears recorded messages.
func (m *MockClient) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishedMsgs = nil
}
