"""Entity type definitions for all MQTT CRDs."""

# MQTTDevice - utility resource for shared device definitions
MQTT_DEVICE = {
    "kind": "MQTTDevice",
    "singular": "mqttdevice",
    "plural": "mqttdevices",
    "short_names": ["dev"],
    "component": None,  # No HA component - utility resource
    "description": "Shared device definition for multiple MQTT entities",
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
    "required": [],
}

# MQTTButton
MQTT_BUTTON = {
    "kind": "MQTTButton",
    "singular": "mqttbutton",
    "plural": "mqttbuttons",
    "short_names": ["btn"],
    "component": "button",
    "description": "Stateless button entity - publishes to command topic when pressed",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish when button is pressed",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "payloadPress": {
            "type": "string",
            "description": "Payload sent when button is pressed (default: PRESS)",
        },
        "deviceClass": {
            "type": "string",
            "description": "Button device class",
            "enum": ["identify", "restart", "update"],
        },
    },
    "required": ["commandTopic"],
}

# MQTTSwitch
MQTT_SWITCH = {
    "kind": "MQTTSwitch",
    "singular": "mqttswitch",
    "plural": "mqttswitches",
    "short_names": [],
    "component": "switch",
    "description": "On/off toggle entity with state feedback",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload representing on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload representing off (default: OFF)",
        },
        "stateOn": {
            "type": "string",
            "description": "State value that means on (if different from payloadOn)",
        },
        "stateOff": {
            "type": "string",
            "description": "State value that means off (if different from payloadOff)",
        },
        "deviceClass": {
            "type": "string",
            "description": "Switch device class",
            "enum": ["outlet", "switch"],
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic"],
}

# MQTTSensor
MQTT_SENSOR = {
    "kind": "MQTTSensor",
    "singular": "mqttsensor",
    "plural": "mqttsensors",
    "short_names": [],
    "component": "sensor",
    "description": "Read-only sensor that reports a value from an MQTT topic",
    "properties": {
        "stateTopic": {
            "type": "string",
            "description": "Topic to read sensor value",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract value from payload",
        },
        "unitOfMeasurement": {
            "type": "string",
            "description": "Unit displayed in HA (e.g. °C, %, W)",
        },
        "deviceClass": {
            "type": "string",
            "description": "HA device class (e.g. temperature, humidity, power, energy, battery)",
        },
        "stateClass": {
            "type": "string",
            "description": "State class for statistics",
            "enum": ["measurement", "total", "total_increasing"],
        },
        "expireAfter": {
            "type": "integer",
            "description": "Seconds after which the sensor value expires",
        },
        "forceUpdate": {
            "type": "boolean",
            "description": "Update HA state even if the value hasn't changed",
        },
        "lastResetValueTemplate": {
            "type": "string",
            "description": "Template for the last reset timestamp",
        },
        "suggestedDisplayPrecision": {
            "type": "integer",
            "description": "Number of decimal places to display",
        },
    },
    "required": ["stateTopic"],
}

# MQTTBinarySensor
MQTT_BINARY_SENSOR = {
    "kind": "MQTTBinarySensor",
    "singular": "mqttbinarysensor",
    "plural": "mqttbinarysensors",
    "short_names": [],
    "component": "binary_sensor",
    "description": "Read-only on/off sensor (e.g. motion detector, door contact)",
    "properties": {
        "stateTopic": {
            "type": "string",
            "description": "Topic to read sensor state",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload representing on/detected (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload representing off/clear (default: OFF)",
        },
        "deviceClass": {
            "type": "string",
            "description": "HA device class (e.g. motion, door, window, moisture, smoke, occupancy)",
        },
        "expireAfter": {
            "type": "integer",
            "description": "Seconds after which the state expires",
        },
        "forceUpdate": {
            "type": "boolean",
            "description": "Update state even if unchanged",
        },
        "offDelay": {
            "type": "integer",
            "description": "Seconds after which the sensor auto-resets to off",
        },
    },
    "required": ["stateTopic"],
}

# MQTTNumber
MQTT_NUMBER = {
    "kind": "MQTTNumber",
    "singular": "mqttnumber",
    "plural": "mqttnumbers",
    "short_names": [],
    "component": "number",
    "description": "Numeric input entity with min/max bounds and step size",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish number value",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current value",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract value from payload",
        },
        "min": {
            "type": "number",
            "description": "Minimum value (default: 1)",
        },
        "max": {
            "type": "number",
            "description": "Maximum value (default: 100)",
        },
        "step": {
            "type": "number",
            "description": "Step size (default: 1)",
        },
        "mode": {
            "type": "string",
            "description": "UI mode",
            "enum": ["auto", "box", "slider"],
        },
        "unitOfMeasurement": {
            "type": "string",
            "description": "Unit displayed in HA",
        },
        "deviceClass": {
            "type": "string",
            "description": "HA device class (e.g. temperature, humidity, power_factor)",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic"],
}

# MQTTSelect
MQTT_SELECT = {
    "kind": "MQTTSelect",
    "singular": "mqttselect",
    "plural": "mqttselects",
    "short_names": [],
    "component": "select",
    "description": "Dropdown selection entity with a fixed list of options",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish selected option",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current selection",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract value from payload",
        },
        "options": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of selectable options",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic", "options"],
}

# MQTTText
MQTT_TEXT = {
    "kind": "MQTTText",
    "singular": "mqtttext",
    "plural": "mqtttexts",
    "short_names": [],
    "component": "text",
    "description": "Free-text input entity",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish text value",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current value",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract value from payload",
        },
        "min": {
            "type": "integer",
            "description": "Minimum text length (default: 0)",
        },
        "max": {
            "type": "integer",
            "description": "Maximum text length (default: 255)",
        },
        "pattern": {
            "type": "string",
            "description": "Regex pattern for validation",
        },
        "mode": {
            "type": "string",
            "description": "Input mode",
            "enum": ["text", "password"],
        },
    },
    "required": ["commandTopic"],
}

# MQTTScene
MQTT_SCENE = {
    "kind": "MQTTScene",
    "singular": "mqttscene",
    "plural": "mqttscenes",
    "short_names": [],
    "component": "scene",
    "description": "Scene entity that can be activated via MQTT",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish when scene is activated",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload sent when scene is activated (default: ON)",
        },
    },
    "required": ["commandTopic"],
}

# MQTTTag
MQTT_TAG = {
    "kind": "MQTTTag",
    "singular": "mqtttag",
    "plural": "mqtttags",
    "short_names": [],
    "component": "tag",
    "description": "Tag scanner entity for NFC, RFID, or QR code scanning",
    "properties": {
        "topic": {
            "type": "string",
            "description": "Topic to subscribe to for tag scans",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract tag ID from payload",
        },
    },
    "required": ["topic"],
}

# MQTTLight
MQTT_LIGHT = {
    "kind": "MQTTLight",
    "singular": "mqttlight",
    "plural": "mqttlights",
    "short_names": [],
    "component": "light",
    "description": "Light entity with optional brightness, color temperature, and RGB color support",
    "properties": {
        # Schema selection
        "schema": {
            "type": "string",
            "description": "Light schema mode",
            "enum": ["default", "json", "template"],
        },
        # Common fields
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current state",
        },
        # Basic schema fields
        "payloadOn": {
            "type": "string",
            "description": "Payload for on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload for off (default: OFF)",
        },
        "brightnessCommandTopic": {
            "type": "string",
            "description": "Topic for brightness commands",
        },
        "brightnessStateTopic": {
            "type": "string",
            "description": "Topic for brightness state",
        },
        "brightnessScale": {
            "type": "integer",
            "description": "Max brightness value (default: 255)",
        },
        "brightnessValueTemplate": {
            "type": "string",
            "description": "Template to extract brightness",
        },
        "colorTempCommandTopic": {
            "type": "string",
            "description": "Topic for color temperature commands",
        },
        "colorTempStateTopic": {
            "type": "string",
            "description": "Topic for color temperature state",
        },
        "colorTempValueTemplate": {
            "type": "string",
            "description": "Template to extract color temp",
        },
        "rgbCommandTopic": {
            "type": "string",
            "description": "Topic for RGB color commands",
        },
        "rgbStateTopic": {
            "type": "string",
            "description": "Topic for RGB color state",
        },
        "rgbCommandTemplate": {
            "type": "string",
            "description": "Template for RGB command payload",
        },
        "rgbValueTemplate": {
            "type": "string",
            "description": "Template to extract RGB state",
        },
        "effectCommandTopic": {
            "type": "string",
            "description": "Topic for effect commands",
        },
        "effectStateTopic": {
            "type": "string",
            "description": "Topic for effect state",
        },
        "effectList": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of supported effects",
        },
        "effectValueTemplate": {
            "type": "string",
            "description": "Template to extract effect",
        },
        "minMireds": {
            "type": "integer",
            "description": "Minimum color temp in mireds",
        },
        "maxMireds": {
            "type": "integer",
            "description": "Maximum color temp in mireds",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
        "onCommandType": {
            "type": "string",
            "description": "On command type",
            "enum": ["last", "first", "brightness"],
        },
        # JSON schema fields
        "brightness": {
            "type": "boolean",
            "description": "Enable brightness support (JSON schema)",
        },
        "colorTemp": {
            "type": "boolean",
            "description": "Enable color temperature (JSON schema)",
        },
        "effect": {
            "type": "boolean",
            "description": "Enable effects (JSON schema)",
        },
        "supportedColorModes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported color modes (e.g. rgb, xy, hs, color_temp)",
        },
        # Template schema fields
        "commandOnTemplate": {
            "type": "string",
            "description": "Template for on command (template schema)",
        },
        "commandOffTemplate": {
            "type": "string",
            "description": "Template for off command (template schema)",
        },
        "stateTemplate": {
            "type": "string",
            "description": "Template to extract state (template schema)",
        },
        "brightnessTemplate": {
            "type": "string",
            "description": "Template to extract brightness (template schema)",
        },
        "colorTempTemplate": {
            "type": "string",
            "description": "Template to extract color temp (template schema)",
        },
        "redTemplate": {
            "type": "string",
            "description": "Template to extract red value (template schema)",
        },
        "greenTemplate": {
            "type": "string",
            "description": "Template to extract green value (template schema)",
        },
        "blueTemplate": {
            "type": "string",
            "description": "Template to extract blue value (template schema)",
        },
    },
    "required": ["commandTopic"],
}

# MQTTCover
MQTT_COVER = {
    "kind": "MQTTCover",
    "singular": "mqttcover",
    "plural": "mqttcovers",
    "short_names": [],
    "component": "cover",
    "description": "Cover entity for garage doors, blinds, shutters, and similar devices",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic for open/close/stop commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read cover state",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "positionTopic": {
            "type": "string",
            "description": "Topic to read current position",
        },
        "setPositionTopic": {
            "type": "string",
            "description": "Topic to publish position commands",
        },
        "setPositionTemplate": {
            "type": "string",
            "description": "Template for position command payload",
        },
        "positionTemplate": {
            "type": "string",
            "description": "Template to extract position from payload",
        },
        "tiltCommandTopic": {
            "type": "string",
            "description": "Topic for tilt commands",
        },
        "tiltStatusTopic": {
            "type": "string",
            "description": "Topic to read tilt position",
        },
        "tiltStatusTemplate": {
            "type": "string",
            "description": "Template to extract tilt from payload",
        },
        "payloadOpen": {
            "type": "string",
            "description": "Payload for open command (default: OPEN)",
        },
        "payloadClose": {
            "type": "string",
            "description": "Payload for close command (default: CLOSE)",
        },
        "payloadStop": {
            "type": "string",
            "description": "Payload for stop command (default: STOP)",
        },
        "stateOpen": {
            "type": "string",
            "description": "State value meaning open (default: open)",
        },
        "stateClosed": {
            "type": "string",
            "description": "State value meaning closed (default: closed)",
        },
        "stateOpening": {
            "type": "string",
            "description": "State value meaning opening (default: opening)",
        },
        "stateClosing": {
            "type": "string",
            "description": "State value meaning closing (default: closing)",
        },
        "stateStopped": {
            "type": "string",
            "description": "State value meaning stopped (default: stopped)",
        },
        "positionOpen": {
            "type": "integer",
            "description": "Position value for fully open (default: 100)",
        },
        "positionClosed": {
            "type": "integer",
            "description": "Position value for fully closed (default: 0)",
        },
        "tiltMin": {
            "type": "integer",
            "description": "Minimum tilt value (default: 0)",
        },
        "tiltMax": {
            "type": "integer",
            "description": "Maximum tilt value (default: 100)",
        },
        "deviceClass": {
            "type": "string",
            "description": "Cover device class",
            "enum": [
                "awning",
                "blind",
                "curtain",
                "damper",
                "door",
                "garage",
                "gate",
                "shade",
                "shutter",
                "window",
            ],
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": [],
}

# MQTTLock
MQTT_LOCK = {
    "kind": "MQTTLock",
    "singular": "mqttlock",
    "plural": "mqttlocks",
    "short_names": [],
    "component": "lock",
    "description": "Lock entity with optional code support",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish lock/unlock commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current lock state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadLock": {
            "type": "string",
            "description": "Payload for lock command (default: LOCK)",
        },
        "payloadUnlock": {
            "type": "string",
            "description": "Payload for unlock command (default: UNLOCK)",
        },
        "payloadOpen": {
            "type": "string",
            "description": "Payload for open command (unlatch)",
        },
        "stateLocked": {
            "type": "string",
            "description": "State value meaning locked (default: LOCKED)",
        },
        "stateUnlocked": {
            "type": "string",
            "description": "State value meaning unlocked (default: UNLOCKED)",
        },
        "stateLocking": {
            "type": "string",
            "description": "State value meaning locking (default: LOCKING)",
        },
        "stateUnlocking": {
            "type": "string",
            "description": "State value meaning unlocking (default: UNLOCKING)",
        },
        "stateJammed": {
            "type": "string",
            "description": "State value meaning jammed (default: JAMMED)",
        },
        "codeFormat": {
            "type": "string",
            "description": "Regex for valid codes (e.g. ^\\d{4}$)",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic"],
}

# MQTTValve
MQTT_VALVE = {
    "kind": "MQTTValve",
    "singular": "mqttvalve",
    "plural": "mqttvalves",
    "short_names": [],
    "component": "valve",
    "description": "Valve entity for controlling water, gas, or irrigation valves",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish open/close commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current valve state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "positionTopic": {
            "type": "string",
            "description": "Topic to read current position",
        },
        "setPositionTopic": {
            "type": "string",
            "description": "Topic to publish position commands",
        },
        "setPositionTemplate": {
            "type": "string",
            "description": "Template for position command payload",
        },
        "positionTemplate": {
            "type": "string",
            "description": "Template to extract position from payload",
        },
        "payloadOpen": {
            "type": "string",
            "description": "Payload for open command (default: OPEN)",
        },
        "payloadClose": {
            "type": "string",
            "description": "Payload for close command (default: CLOSE)",
        },
        "payloadStop": {
            "type": "string",
            "description": "Payload for stop command (default: STOP)",
        },
        "stateOpen": {
            "type": "string",
            "description": "State value meaning open (default: open)",
        },
        "stateClosed": {
            "type": "string",
            "description": "State value meaning closed (default: closed)",
        },
        "stateOpening": {
            "type": "string",
            "description": "State value meaning opening (default: opening)",
        },
        "stateClosing": {
            "type": "string",
            "description": "State value meaning closing (default: closing)",
        },
        "deviceClass": {
            "type": "string",
            "description": "Valve device class",
            "enum": ["water", "gas"],
        },
        "reportsPosition": {
            "type": "boolean",
            "description": "Whether the valve reports position",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": [],
}

# MQTTFan
MQTT_FAN = {
    "kind": "MQTTFan",
    "singular": "mqttfan",
    "plural": "mqttfans",
    "short_names": [],
    "component": "fan",
    "description": "Fan entity with speed, direction, oscillation, and preset mode support",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current on/off state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload for on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload for off (default: OFF)",
        },
        "percentageCommandTopic": {
            "type": "string",
            "description": "Topic for speed percentage commands",
        },
        "percentageStateTopic": {
            "type": "string",
            "description": "Topic to read speed percentage",
        },
        "percentageCommandTemplate": {
            "type": "string",
            "description": "Template for percentage command",
        },
        "percentageValueTemplate": {
            "type": "string",
            "description": "Template to extract percentage",
        },
        "speedRangeMin": {
            "type": "integer",
            "description": "Minimum speed value (default: 1)",
        },
        "speedRangeMax": {
            "type": "integer",
            "description": "Maximum speed value (default: 100)",
        },
        "presetModeCommandTopic": {
            "type": "string",
            "description": "Topic for preset mode commands",
        },
        "presetModeStateTopic": {
            "type": "string",
            "description": "Topic to read preset mode",
        },
        "presetModeCommandTemplate": {
            "type": "string",
            "description": "Template for preset mode command",
        },
        "presetModeValueTemplate": {
            "type": "string",
            "description": "Template to extract preset mode",
        },
        "presetModes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of supported preset modes",
        },
        "oscillationCommandTopic": {
            "type": "string",
            "description": "Topic for oscillation commands",
        },
        "oscillationStateTopic": {
            "type": "string",
            "description": "Topic to read oscillation state",
        },
        "oscillationCommandTemplate": {
            "type": "string",
            "description": "Template for oscillation command",
        },
        "oscillationValueTemplate": {
            "type": "string",
            "description": "Template to extract oscillation state",
        },
        "payloadOscillationOn": {
            "type": "string",
            "description": "Payload for oscillation on (default: oscillate_on)",
        },
        "payloadOscillationOff": {
            "type": "string",
            "description": "Payload for oscillation off (default: oscillate_off)",
        },
        "directionCommandTopic": {
            "type": "string",
            "description": "Topic for direction commands",
        },
        "directionStateTopic": {
            "type": "string",
            "description": "Topic to read direction state",
        },
        "directionValueTemplate": {
            "type": "string",
            "description": "Template to extract direction",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic"],
}

# MQTTSiren
MQTT_SIREN = {
    "kind": "MQTTSiren",
    "singular": "mqttsiren",
    "plural": "mqttsirens",
    "short_names": [],
    "component": "siren",
    "description": "Siren entity with optional tone, volume, and duration support",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload for on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload for off (default: OFF)",
        },
        "stateOn": {
            "type": "string",
            "description": "State value meaning on (default: ON)",
        },
        "stateOff": {
            "type": "string",
            "description": "State value meaning off (default: OFF)",
        },
        "availableTones": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of supported tones",
        },
        "supportTurnOn": {
            "type": "boolean",
            "description": "Whether the siren supports turn on (default: true)",
        },
        "supportTurnOff": {
            "type": "boolean",
            "description": "Whether the siren supports turn off (default: true)",
        },
        "supportDuration": {
            "type": "boolean",
            "description": "Whether duration is supported (default: true)",
        },
        "supportVolumeSet": {
            "type": "boolean",
            "description": "Whether volume level is supported (default: true)",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic"],
}

# MQTTCamera
MQTT_CAMERA = {
    "kind": "MQTTCamera",
    "singular": "mqttcamera",
    "plural": "mqttcameras",
    "short_names": [],
    "component": "camera",
    "description": "Camera entity that receives images via MQTT",
    "properties": {
        "topic": {
            "type": "string",
            "description": "MQTT topic to subscribe to for image data",
        },
        "imageEncoding": {
            "type": "string",
            "description": "Image encoding (b64 for base64-encoded images)",
        },
        "stateClass": {
            "type": "string",
            "description": "State class for statistics",
            "enum": ["measurement", "total", "total_increasing"],
        },
        "expireAfter": {
            "type": "integer",
            "description": "Seconds after which the image expires",
            "minimum": 0,
        },
    },
    "required": ["topic"],
}

# MQTTImage
MQTT_IMAGE = {
    "kind": "MQTTImage",
    "singular": "mqttimage",
    "plural": "mqttimages",
    "short_names": [],
    "component": "image",
    "description": "Image entity that displays a static image from an MQTT topic or URL",
    "properties": {
        "imageTopic": {
            "type": "string",
            "description": "Topic to receive raw image data",
        },
        "imageEncoding": {
            "type": "string",
            "description": "Image encoding (b64 for base64-encoded images)",
        },
        "urlTopic": {
            "type": "string",
            "description": "Topic to receive image URL",
        },
        "urlTemplate": {
            "type": "string",
            "description": "Template to extract URL from payload",
        },
        "contentType": {
            "type": "string",
            "description": "Image MIME type (default: image/png)",
        },
    },
    "required": [],
}

# MQTTNotify
MQTT_NOTIFY = {
    "kind": "MQTTNotify",
    "singular": "mqttnotify",
    "plural": "mqttnotifys",
    "short_names": [],
    "component": "notify",
    "description": "Notification service entity that sends messages to a device via MQTT",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish notification messages",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the notification payload",
        },
    },
    "required": ["commandTopic"],
}

# MQTTUpdate
MQTT_UPDATE = {
    "kind": "MQTTUpdate",
    "singular": "mqttupdate",
    "plural": "mqttupdates",
    "short_names": [],
    "component": "update",
    "description": "Firmware/software update entity that tracks available updates via MQTT",
    "properties": {
        "stateTopic": {
            "type": "string",
            "description": "Topic with JSON payload containing update info",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "commandTopic": {
            "type": "string",
            "description": "Topic to trigger update installation",
        },
        "payloadInstall": {
            "type": "string",
            "description": "Payload to trigger installation (default: INSTALL)",
        },
        "latestVersionTopic": {
            "type": "string",
            "description": "Topic to read latest available version",
        },
        "latestVersionTemplate": {
            "type": "string",
            "description": "Template to extract latest version",
        },
        "deviceClass": {
            "type": "string",
            "description": "Update device class",
            "enum": ["firmware"],
        },
        "entityPicture": {
            "type": "string",
            "description": "URL to an image for the update entity",
        },
        "releaseUrl": {
            "type": "string",
            "description": "URL to release notes",
        },
        "releaseSummary": {
            "type": "string",
            "description": "Summary of the release",
        },
        "title": {
            "type": "string",
            "description": "Title of the software/firmware",
        },
    },
    "required": ["stateTopic"],
}

# MQTTClimate
MQTT_CLIMATE = {
    "kind": "MQTTClimate",
    "singular": "mqttclimate",
    "plural": "mqttclimates",
    "short_names": ["hvac"],
    "component": "climate",
    "description": "Thermostat/HVAC entity with temperature control, modes, and fan speed",
    "properties": {
        "temperatureCommandTopic": {
            "type": "string",
            "description": "Topic to set target temperature",
        },
        "temperatureStateTopic": {
            "type": "string",
            "description": "Topic to read target temperature",
        },
        "temperatureCommandTemplate": {
            "type": "string",
            "description": "Template for temperature command",
        },
        "temperatureStateTemplate": {
            "type": "string",
            "description": "Template to extract target temp",
        },
        "currentTemperatureTopic": {
            "type": "string",
            "description": "Topic to read current temperature",
        },
        "currentTemperatureTemplate": {
            "type": "string",
            "description": "Template to extract current temp",
        },
        "modeCommandTopic": {
            "type": "string",
            "description": "Topic to set HVAC mode",
        },
        "modeStateTopic": {
            "type": "string",
            "description": "Topic to read HVAC mode",
        },
        "modeCommandTemplate": {
            "type": "string",
            "description": "Template for mode command",
        },
        "modeStateTemplate": {
            "type": "string",
            "description": "Template to extract mode",
        },
        "modes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported HVAC modes",
        },
        "fanModeCommandTopic": {
            "type": "string",
            "description": "Topic to set fan mode",
        },
        "fanModeStateTopic": {
            "type": "string",
            "description": "Topic to read fan mode",
        },
        "fanModeCommandTemplate": {
            "type": "string",
            "description": "Template for fan mode command",
        },
        "fanModeStateTemplate": {
            "type": "string",
            "description": "Template to extract fan mode",
        },
        "fanModes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported fan modes",
        },
        "swingModeCommandTopic": {
            "type": "string",
            "description": "Topic to set swing mode",
        },
        "swingModeStateTopic": {
            "type": "string",
            "description": "Topic to read swing mode",
        },
        "swingModes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported swing modes",
        },
        "presetModeCommandTopic": {
            "type": "string",
            "description": "Topic to set preset mode",
        },
        "presetModeStateTopic": {
            "type": "string",
            "description": "Topic to read preset mode",
        },
        "presetModes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported preset modes (e.g. away, eco, boost)",
        },
        "actionTopic": {
            "type": "string",
            "description": "Topic to read current HVAC action",
        },
        "actionTemplate": {
            "type": "string",
            "description": "Template to extract action",
        },
        "tempStep": {
            "type": "number",
            "description": "Step size for temperature adjustments (default: 1)",
        },
        "minTemp": {
            "type": "number",
            "description": "Minimum setpoint temperature",
        },
        "maxTemp": {
            "type": "number",
            "description": "Maximum setpoint temperature",
        },
        "temperatureUnit": {
            "type": "string",
            "description": "Temperature unit",
            "enum": ["C", "F"],
        },
        "precision": {
            "type": "number",
            "description": "Temperature precision (default: 0.1)",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": [],
}

# MQTTHumidifier
MQTT_HUMIDIFIER = {
    "kind": "MQTTHumidifier",
    "singular": "mqtthumidifier",
    "plural": "mqtthumidifiers",
    "short_names": [],
    "component": "humidifier",
    "description": "Humidifier entity with target humidity and mode support",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read current on/off state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload for on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload for off (default: OFF)",
        },
        "targetHumidityCommandTopic": {
            "type": "string",
            "description": "Topic to set target humidity",
        },
        "targetHumidityStateTopic": {
            "type": "string",
            "description": "Topic to read target humidity",
        },
        "targetHumidityCommandTemplate": {
            "type": "string",
            "description": "Template for target humidity command",
        },
        "targetHumidityStateTemplate": {
            "type": "string",
            "description": "Template to extract target humidity",
        },
        "currentHumidityTopic": {
            "type": "string",
            "description": "Topic to read current humidity",
        },
        "currentHumidityTemplate": {
            "type": "string",
            "description": "Template to extract current humidity",
        },
        "modeCommandTopic": {
            "type": "string",
            "description": "Topic to set mode",
        },
        "modeStateTopic": {
            "type": "string",
            "description": "Topic to read current mode",
        },
        "modeCommandTemplate": {
            "type": "string",
            "description": "Template for mode command",
        },
        "modeStateTemplate": {
            "type": "string",
            "description": "Template to extract mode",
        },
        "modes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported modes (e.g. normal, eco, boost, sleep)",
        },
        "actionTopic": {
            "type": "string",
            "description": "Topic to read current action",
        },
        "actionTemplate": {
            "type": "string",
            "description": "Template to extract action",
        },
        "minHumidity": {
            "type": "number",
            "description": "Minimum target humidity (default: 0)",
        },
        "maxHumidity": {
            "type": "number",
            "description": "Maximum target humidity (default: 100)",
        },
        "deviceClass": {
            "type": "string",
            "description": "Humidifier device class",
            "enum": ["humidifier", "dehumidifier"],
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": ["commandTopic", "targetHumidityCommandTopic"],
}

# MQTTWaterHeater
MQTT_WATER_HEATER = {
    "kind": "MQTTWaterHeater",
    "singular": "mqttwaterheater",
    "plural": "mqttwaterheaters",
    "short_names": [],
    "component": "water_heater",
    "description": "Water heater entity with temperature control and operation modes",
    "properties": {
        "temperatureCommandTopic": {
            "type": "string",
            "description": "Topic to set target temperature",
        },
        "temperatureStateTopic": {
            "type": "string",
            "description": "Topic to read target temperature",
        },
        "temperatureCommandTemplate": {
            "type": "string",
            "description": "Template for temperature command",
        },
        "temperatureStateTemplate": {
            "type": "string",
            "description": "Template to extract target temp",
        },
        "currentTemperatureTopic": {
            "type": "string",
            "description": "Topic to read current temperature",
        },
        "currentTemperatureTemplate": {
            "type": "string",
            "description": "Template to extract current temp",
        },
        "modeCommandTopic": {
            "type": "string",
            "description": "Topic to set operation mode",
        },
        "modeStateTopic": {
            "type": "string",
            "description": "Topic to read operation mode",
        },
        "modeCommandTemplate": {
            "type": "string",
            "description": "Template for mode command",
        },
        "modeStateTemplate": {
            "type": "string",
            "description": "Template to extract mode",
        },
        "modes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported modes (e.g. off, eco, electric, gas, heat_pump, high_demand, performance)",
        },
        "powerCommandTopic": {
            "type": "string",
            "description": "Topic to publish on/off commands",
        },
        "payloadOn": {
            "type": "string",
            "description": "Payload for on (default: ON)",
        },
        "payloadOff": {
            "type": "string",
            "description": "Payload for off (default: OFF)",
        },
        "minTemp": {
            "type": "number",
            "description": "Minimum target temperature (default: 110)",
        },
        "maxTemp": {
            "type": "number",
            "description": "Maximum target temperature (default: 140)",
        },
        "temperatureUnit": {
            "type": "string",
            "description": "Temperature unit",
            "enum": ["C", "F"],
        },
        "precision": {
            "type": "number",
            "description": "Temperature precision (default: 0.1)",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": [],
}

# MQTTVacuum
MQTT_VACUUM = {
    "kind": "MQTTVacuum",
    "singular": "mqttvacuum",
    "plural": "mqttvacuums",
    "short_names": [],
    "component": "vacuum",
    "description": "Robot vacuum entity with start, stop, pause, return to base, and cleaning features",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic for basic commands (start, stop, return_to_base, etc.)",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read vacuum state",
        },
        "sendCommandTopic": {
            "type": "string",
            "description": "Topic for custom commands",
        },
        "setFanSpeedTopic": {
            "type": "string",
            "description": "Topic for fan speed commands",
        },
        "fanSpeedList": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of supported fan speeds",
        },
        "payloadStart": {
            "type": "string",
            "description": "Payload for start command (default: start)",
        },
        "payloadStop": {
            "type": "string",
            "description": "Payload for stop command (default: stop)",
        },
        "payloadPause": {
            "type": "string",
            "description": "Payload for pause command (default: pause)",
        },
        "payloadReturnToBase": {
            "type": "string",
            "description": "Payload for return to base command (default: return_to_base)",
        },
        "payloadCleanSpot": {
            "type": "string",
            "description": "Payload for clean spot command (default: clean_spot)",
        },
        "payloadLocate": {
            "type": "string",
            "description": "Payload for locate command (default: locate)",
        },
        "supportedFeatures": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported features (e.g. start, stop, pause, return_home, fan_speed, send_command, locate, clean_spot)",
        },
        "schema": {
            "type": "string",
            "description": "Vacuum schema",
            "enum": ["legacy", "state"],
        },
    },
    "required": [],
}

# MQTTLawnMower
MQTT_LAWN_MOWER = {
    "kind": "MQTTLawnMower",
    "singular": "mqttlawnmower",
    "plural": "mqttlawnmowers",
    "short_names": [],
    "component": "lawn_mower",
    "description": "Robot lawn mower entity with start mowing, pause, and dock commands",
    "properties": {
        "activityStateTopic": {
            "type": "string",
            "description": "Topic to read mower activity state",
        },
        "activityValueTemplate": {
            "type": "string",
            "description": "Template to extract activity from payload",
        },
        "dockCommandTopic": {
            "type": "string",
            "description": "Topic to publish dock command",
        },
        "dockCommandTemplate": {
            "type": "string",
            "description": "Template for dock command payload",
        },
        "pauseCommandTopic": {
            "type": "string",
            "description": "Topic to publish pause command",
        },
        "pauseCommandTemplate": {
            "type": "string",
            "description": "Template for pause command payload",
        },
        "startMowingCommandTopic": {
            "type": "string",
            "description": "Topic to publish start mowing command",
        },
        "startMowingCommandTemplate": {
            "type": "string",
            "description": "Template for start mowing command payload",
        },
        "optimistic": {
            "type": "boolean",
            "description": "Assume state changes immediately",
        },
    },
    "required": [],
}

# MQTTAlarmControlPanel
MQTT_ALARM_CONTROL_PANEL = {
    "kind": "MQTTAlarmControlPanel",
    "singular": "mqttalarmcontrolpanel",
    "plural": "mqttalarmcontrolpanels",
    "short_names": [],
    "component": "alarm_control_panel",
    "description": "Alarm control panel entity with arm/disarm modes and optional code support",
    "properties": {
        "commandTopic": {
            "type": "string",
            "description": "Topic to publish arm/disarm commands",
        },
        "stateTopic": {
            "type": "string",
            "description": "Topic to read alarm state",
        },
        "commandTemplate": {
            "type": "string",
            "description": "Template for the command payload",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadArmHome": {
            "type": "string",
            "description": "Payload for arm home (default: ARM_HOME)",
        },
        "payloadArmAway": {
            "type": "string",
            "description": "Payload for arm away (default: ARM_AWAY)",
        },
        "payloadArmNight": {
            "type": "string",
            "description": "Payload for arm night (default: ARM_NIGHT)",
        },
        "payloadArmVacation": {
            "type": "string",
            "description": "Payload for arm vacation (default: ARM_VACATION)",
        },
        "payloadArmCustomBypass": {
            "type": "string",
            "description": "Payload for arm custom bypass (default: ARM_CUSTOM_BYPASS)",
        },
        "payloadDisarm": {
            "type": "string",
            "description": "Payload for disarm (default: DISARM)",
        },
        "payloadTrigger": {
            "type": "string",
            "description": "Payload for trigger",
        },
        "codeArmRequired": {
            "type": "boolean",
            "description": "Whether code is required to arm (default: true)",
        },
        "codeDisarmRequired": {
            "type": "boolean",
            "description": "Whether code is required to disarm (default: true)",
        },
        "codeTriggerRequired": {
            "type": "boolean",
            "description": "Whether code is required to trigger (default: true)",
        },
        "codeFormat": {
            "type": "string",
            "description": "Code format",
            "enum": ["number", "text"],
        },
        "supportedFeatures": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Supported features (e.g. arm_home, arm_away, arm_night, trigger)",
        },
    },
    "required": ["commandTopic", "stateTopic"],
}

# MQTTDeviceTracker
MQTT_DEVICE_TRACKER = {
    "kind": "MQTTDeviceTracker",
    "singular": "mqttdevicetracker",
    "plural": "mqttdevicetrackers",
    "short_names": [],
    "component": "device_tracker",
    "description": "Device tracker entity for presence detection and location tracking via MQTT",
    "properties": {
        "stateTopic": {
            "type": "string",
            "description": "Topic to read tracker state (home/not_home or zone name)",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract state from payload",
        },
        "payloadHome": {
            "type": "string",
            "description": "Payload representing home (default: home)",
        },
        "payloadNotHome": {
            "type": "string",
            "description": "Payload representing not home (default: not_home)",
        },
        "payloadReset": {
            "type": "string",
            "description": "Payload that resets the tracker to unknown",
        },
        "sourceType": {
            "type": "string",
            "description": "Source type (e.g. gps, router, bluetooth, bluetooth_le)",
        },
    },
    "required": ["stateTopic"],
}

# MQTTDeviceTrigger
MQTT_DEVICE_TRIGGER = {
    "kind": "MQTTDeviceTrigger",
    "singular": "mqttdevicetrigger",
    "plural": "mqttdevicetriggers",
    "short_names": [],
    "component": "device_automation",
    "description": "Device automation trigger that fires when a specific MQTT message is received",
    "properties": {
        "topic": {
            "type": "string",
            "description": "MQTT topic to subscribe to for trigger events",
        },
        "payload": {
            "type": "string",
            "description": "Specific payload that triggers the automation",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract value from payload",
        },
        "type": {
            "type": "string",
            "description": "Trigger type (e.g. button_short_press, button_long_press)",
        },
        "subtype": {
            "type": "string",
            "description": "Trigger subtype (e.g. button_1, turn_on)",
        },
        "automationType": {
            "type": "string",
            "description": "Automation type (always trigger)",
        },
    },
    "required": ["topic", "type", "subtype"],
}

# MQTTEvent
MQTT_EVENT = {
    "kind": "MQTTEvent",
    "singular": "mqttevent",
    "plural": "mqttevents",
    "short_names": [],
    "component": "event",
    "description": "Event entity for stateless events such as button presses or doorbell rings",
    "properties": {
        "stateTopic": {
            "type": "string",
            "description": "Topic to subscribe to for events",
        },
        "valueTemplate": {
            "type": "string",
            "description": "Template to extract event type from payload",
        },
        "eventTypes": {
            "type": "array",
            "items": {"type": "string"},
            "description": "List of supported event types",
        },
        "deviceClass": {
            "type": "string",
            "description": "Event device class",
            "enum": ["button", "doorbell", "motion"],
        },
    },
    "required": ["stateTopic", "eventTypes"],
}


# All entities in order of complexity (for phased generation)
ALL_ENTITIES = [
    # Phase 1: Simplest
    MQTT_DEVICE,
    MQTT_BUTTON,
    # Phase 2: Simple types
    MQTT_SWITCH,
    MQTT_SENSOR,
    MQTT_BINARY_SENSOR,
    MQTT_NUMBER,
    MQTT_SELECT,
    MQTT_TEXT,
    MQTT_SCENE,
    MQTT_TAG,
    # Phase 3: Medium complexity
    MQTT_LIGHT,
    MQTT_COVER,
    MQTT_LOCK,
    MQTT_VALVE,
    MQTT_FAN,
    MQTT_SIREN,
    MQTT_CAMERA,
    MQTT_IMAGE,
    MQTT_NOTIFY,
    MQTT_UPDATE,
    # Phase 4: Complex types
    MQTT_CLIMATE,
    MQTT_HUMIDIFIER,
    MQTT_WATER_HEATER,
    MQTT_VACUUM,
    MQTT_LAWN_MOWER,
    MQTT_ALARM_CONTROL_PANEL,
    MQTT_DEVICE_TRACKER,
    MQTT_DEVICE_TRIGGER,
    MQTT_EVENT,
]
