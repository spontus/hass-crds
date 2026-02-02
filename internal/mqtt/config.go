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
	"fmt"
	"os"
	"strconv"
)

// Config holds the MQTT connection configuration.
type Config struct {
	Broker   string
	Port     int
	ClientID string
	Username string
	Password string
	UseTLS   bool
}

// NewConfigFromEnv creates a Config from environment variables.
func NewConfigFromEnv() (*Config, error) {
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		return nil, fmt.Errorf("MQTT_BROKER environment variable is required")
	}

	portStr := os.Getenv("MQTT_PORT")
	port := 1883
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MQTT_PORT: %w", err)
		}
	}

	clientID := os.Getenv("MQTT_CLIENT_ID")
	if clientID == "" {
		clientID = "hass-crds-controller"
	}

	useTLS := false
	if tlsStr := os.Getenv("MQTT_USE_TLS"); tlsStr == "true" || tlsStr == "1" {
		useTLS = true
	}

	return &Config{
		Broker:   broker,
		Port:     port,
		ClientID: clientID,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
		UseTLS:   useTLS,
	}, nil
}

// BrokerURL returns the full broker URL.
func (c *Config) BrokerURL() string {
	scheme := "tcp"
	if c.UseTLS {
		scheme = "ssl"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, c.Broker, c.Port)
}
