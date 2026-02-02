"""Common schema definitions shared across all MQTT entity CRDs."""

# Entity metadata fields (common to all entity types)
ENTITY_METADATA = {
    "name": {
        "type": "string",
        "description": "Display name in Home Assistant",
    },
    "uniqueId": {
        "type": "string",
        "description": "Unique identifier for HA entity registry (defaults to <namespace>-<name>)",
    },
    "icon": {
        "type": "string",
        "description": "MDI icon (e.g. mdi:thermometer)",
    },
    "entityCategory": {
        "type": "string",
        "description": "Entity category",
        "enum": ["config", "diagnostic"],
    },
    "enabledByDefault": {
        "type": "boolean",
        "description": "Whether the entity is enabled when first discovered",
    },
    "objectId": {
        "type": "string",
        "description": "Override for HA entity ID generation",
    },
}

# Device block (inline device definition)
DEVICE_BLOCK = {
    "device": {
        "type": "object",
        "description": "Device configuration for Home Assistant device registry",
        "properties": {
            "name": {
                "type": "string",
                "description": "Device display name",
            },
            "identifiers": {
                "type": "array",
                "items": {"type": "string"},
                "description": "List of identifiers (at least one of identifiers or connections is needed)",
            },
            "connections": {
                "type": "array",
                "items": {
                    "type": "array",
                    "items": {"type": "string"},
                },
                "description": "List of [type, value] pairs (e.g. [[mac, aa:bb:cc:dd:ee:ff]])",
            },
            "manufacturer": {
                "type": "string",
                "description": "Device manufacturer",
            },
            "model": {
                "type": "string",
                "description": "Device model",
            },
            "modelId": {
                "type": "string",
                "description": "Device model identifier",
            },
            "serialNumber": {
                "type": "string",
                "description": "Device serial number",
            },
            "hwVersion": {
                "type": "string",
                "description": "Hardware version",
            },
            "swVersion": {
                "type": "string",
                "description": "Software version",
            },
            "suggestedArea": {
                "type": "string",
                "description": "Suggested area in Home Assistant (e.g. Living Room)",
            },
            "configurationUrl": {
                "type": "string",
                "description": "URL for device configuration",
            },
            "viaDevice": {
                "type": "string",
                "description": "Identifier of device that routes messages",
            },
        },
    },
}

# Device reference (reference to MQTTDevice resource)
DEVICE_REF = {
    "deviceRef": {
        "type": "object",
        "description": "Reference to an MQTTDevice resource instead of inline device block",
        "properties": {
            "name": {
                "type": "string",
                "description": "Name of an MQTTDevice resource in the same namespace",
            },
        },
        "required": ["name"],
    },
}

# Availability configuration
AVAILABILITY = {
    "availability": {
        "type": "array",
        "description": "List of availability topics",
        "items": {
            "type": "object",
            "properties": {
                "topic": {
                    "type": "string",
                    "description": "MQTT topic for availability",
                },
                "payloadAvailable": {
                    "type": "string",
                    "description": "Payload indicating available (default: online)",
                },
                "payloadNotAvailable": {
                    "type": "string",
                    "description": "Payload indicating unavailable (default: offline)",
                },
                "valueTemplate": {
                    "type": "string",
                    "description": "Template to extract availability from payload",
                },
            },
            "required": ["topic"],
        },
    },
    "availabilityMode": {
        "type": "string",
        "description": "How to combine multiple availability topics",
        "enum": ["all", "any", "latest"],
    },
}

# MQTT options
MQTT_OPTIONS = {
    "qos": {
        "type": "integer",
        "description": "MQTT QoS level",
        "minimum": 0,
        "maximum": 2,
    },
    "retain": {
        "type": "boolean",
        "description": "Whether to retain messages on command/state topics",
    },
    "encoding": {
        "type": "string",
        "description": "Payload encoding (default: utf-8)",
    },
}

# JSON attributes
JSON_ATTRIBUTES = {
    "jsonAttributesTopic": {
        "type": "string",
        "description": "MQTT topic for JSON attributes",
    },
    "jsonAttributesTemplate": {
        "type": "string",
        "description": "Template to extract attributes from payload",
    },
}

# Rediscovery interval (controller-only field)
REDISCOVERY = {
    "rediscoverInterval": {
        "type": "string",
        "description": "How often to re-publish the discovery config payload (e.g. 5m, 1h)",
    },
}

# Status subresource schema
STATUS_SCHEMA = {
    "type": "object",
    "properties": {
        "lastPublished": {
            "type": "string",
            "format": "date-time",
            "description": "Timestamp of last discovery publish",
        },
        "discoveryTopic": {
            "type": "string",
            "description": "MQTT discovery topic path",
        },
        "conditions": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "type": {
                        "type": "string",
                        "description": "Condition type (Published, MQTTConnected)",
                    },
                    "status": {
                        "type": "string",
                        "enum": ["True", "False", "Unknown"],
                    },
                    "lastTransitionTime": {
                        "type": "string",
                        "format": "date-time",
                    },
                    "reason": {
                        "type": "string",
                    },
                    "message": {
                        "type": "string",
                    },
                },
                "required": ["type", "status"],
            },
        },
    },
}


def get_all_common_properties() -> dict:
    """Returns all common properties merged together."""
    props = {}
    props.update(ENTITY_METADATA)
    props.update(DEVICE_BLOCK)
    props.update(DEVICE_REF)
    props.update(AVAILABILITY)
    props.update(MQTT_OPTIONS)
    props.update(JSON_ATTRIBUTES)
    props.update(REDISCOVERY)
    return props
