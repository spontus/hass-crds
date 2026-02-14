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

package payload

import (
	"encoding/json"
	"strings"
	"unicode"
)

const (
	// OriginName is the name used in the origin block of discovery payloads.
	OriginName = "hass-crds"

	// OriginSupportURL is the support URL used in the origin block.
	OriginSupportURL = "https://github.com/spontus/hass-crds"
)

// DefaultOrigin returns the default origin block for discovery payloads.
func DefaultOrigin() map[string]interface{} {
	return map[string]interface{}{
		"name":        OriginName,
		"support_url": OriginSupportURL,
	}
}

// Builder builds MQTT discovery payloads.
type Builder struct {
	data map[string]interface{}
}

// New creates a new payload builder.
func New() *Builder {
	return &Builder{
		data: make(map[string]interface{}),
	}
}

// Set adds a key-value pair to the payload.
// The key is converted from camelCase to snake_case for Home Assistant compatibility.
func (b *Builder) Set(key string, value interface{}) *Builder {
	if value == nil {
		return b
	}

	// Check for zero values that shouldn't be included
	switch v := value.(type) {
	case string:
		if v == "" {
			return b
		}
	case []string:
		if len(v) == 0 {
			return b
		}
	case [][]string:
		if len(v) == 0 {
			return b
		}
	case *bool:
		if v == nil {
			return b
		}
		b.data[camelToSnake(key)] = *v
		return b
	case *int:
		if v == nil {
			return b
		}
		b.data[camelToSnake(key)] = *v
		return b
	case *float64:
		if v == nil {
			return b
		}
		b.data[camelToSnake(key)] = *v
		return b
	}

	b.data[camelToSnake(key)] = value
	return b
}

// SetRaw adds a key-value pair without converting the key.
func (b *Builder) SetRaw(key string, value interface{}) *Builder {
	if value == nil {
		return b
	}
	b.data[key] = value
	return b
}

// SetDevice adds a device block to the payload.
func (b *Builder) SetDevice(device map[string]interface{}) *Builder {
	if len(device) > 0 {
		b.data["device"] = device
	}
	return b
}

// SetAvailability adds availability configuration to the payload.
func (b *Builder) SetAvailability(availability []map[string]interface{}) *Builder {
	if len(availability) > 0 {
		b.data["availability"] = availability
	}
	return b
}

// SetOrigin adds an origin block to the payload.
func (b *Builder) SetOrigin(origin map[string]interface{}) *Builder {
	if len(origin) > 0 {
		b.data["origin"] = origin
	}
	return b
}

// Build returns the payload as JSON bytes.
func (b *Builder) Build() ([]byte, error) {
	return json.Marshal(b.data)
}

// BuildMap returns the payload as a map.
func (b *Builder) BuildMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range b.data {
		result[k] = v
	}
	return result
}

// camelToSnake converts a camelCase string to snake_case.
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// DeviceBlockToMap converts a device spec to a map suitable for the payload.
func DeviceBlockToMap(name string, identifiers []string, connections [][]string,
	manufacturer, model, modelId, serialNumber, hwVersion, swVersion,
	suggestedArea, configurationUrl, viaDevice string) map[string]interface{} {

	device := make(map[string]interface{})

	if name != "" {
		device["name"] = name
	}
	if len(identifiers) > 0 {
		device["identifiers"] = identifiers
	}
	if len(connections) > 0 {
		device["connections"] = connections
	}
	if manufacturer != "" {
		device["manufacturer"] = manufacturer
	}
	if model != "" {
		device["model"] = model
	}
	if modelId != "" {
		device["model_id"] = modelId
	}
	if serialNumber != "" {
		device["serial_number"] = serialNumber
	}
	if hwVersion != "" {
		device["hw_version"] = hwVersion
	}
	if swVersion != "" {
		device["sw_version"] = swVersion
	}
	if suggestedArea != "" {
		device["suggested_area"] = suggestedArea
	}
	if configurationUrl != "" {
		device["configuration_url"] = configurationUrl
	}
	if viaDevice != "" {
		device["via_device"] = viaDevice
	}

	return device
}

// AvailabilityToMap converts availability configs to maps suitable for the payload.
func AvailabilityToMap(topic, payloadAvailable, payloadNotAvailable, valueTemplate string) map[string]interface{} {
	avail := make(map[string]interface{})
	avail["topic"] = topic

	if payloadAvailable != "" {
		avail["payload_available"] = payloadAvailable
	}
	if payloadNotAvailable != "" {
		avail["payload_not_available"] = payloadNotAvailable
	}
	if valueTemplate != "" {
		avail["value_template"] = valueTemplate
	}

	return avail
}
