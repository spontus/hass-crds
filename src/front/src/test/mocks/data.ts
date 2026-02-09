import type {
  EntityType,
  EntitySummary,
  NamespaceSummary,
  SchemaProperty,
} from '../../types/entity'

export const mockEntityTypes: EntityType[] = [
  { kind: 'MQTTButton', plural: 'mqttbuttons', description: 'Stateless button', category: 'Controls' },
  { kind: 'MQTTSwitch', plural: 'mqttswitches', description: 'On/off switch', category: 'Controls' },
  { kind: 'MQTTSensor', plural: 'mqttsensors', description: 'Read-only sensor', category: 'Sensors' },
  { kind: 'MQTTBinarySensor', plural: 'mqttbinarysensors', description: 'Binary sensor', category: 'Sensors' },
  { kind: 'MQTTLight', plural: 'mqttlights', description: 'Light control', category: 'Lighting' },
]

export const mockCategories: Record<string, EntityType[]> = {
  Controls: [mockEntityTypes[0], mockEntityTypes[1]],
  Sensors: [mockEntityTypes[2], mockEntityTypes[3]],
  Lighting: [mockEntityTypes[4]],
}

export const mockEntities: EntitySummary[] = [
  {
    kind: 'MQTTButton',
    apiVersion: 'mqtt.home-assistant.io/v1alpha1',
    name: 'test-button',
    namespace: 'default',
    displayName: 'Test Button',
    published: true,
    createdAt: '2024-01-15T10:00:00Z',
    labels: {},
  },
  {
    kind: 'MQTTSensor',
    apiVersion: 'mqtt.home-assistant.io/v1alpha1',
    name: 'temp-sensor',
    namespace: 'default',
    displayName: 'Temperature Sensor',
    published: true,
    createdAt: '2024-01-14T09:00:00Z',
    labels: {},
  },
  {
    kind: 'MQTTSwitch',
    apiVersion: 'mqtt.home-assistant.io/v1alpha1',
    name: 'power-switch',
    namespace: 'production',
    displayName: 'Power Switch',
    published: false,
    createdAt: '2024-01-13T08:00:00Z',
    labels: {},
  },
]

export const mockNamespaces: NamespaceSummary[] = [
  { name: 'default', status: 'Active', labels: {} },
  { name: 'production', status: 'Active', labels: { env: 'prod' } },
  { name: 'staging', status: 'Active', labels: { env: 'staging' } },
]

export const mockButtonSchema: SchemaProperty = {
  type: 'object',
  required: ['commandTopic'],
  properties: {
    name: {
      type: 'string',
      description: 'Display name in Home Assistant',
    },
    commandTopic: {
      type: 'string',
      description: 'MQTT topic to publish on press',
    },
    payloadPress: {
      type: 'string',
      description: 'Payload to send on press',
      default: 'PRESS',
    },
    deviceClass: {
      type: 'string',
      description: 'Device class for the button',
      enum: ['identify', 'restart', 'update'],
    },
    icon: {
      type: 'string',
      description: 'Material Design icon',
    },
    uniqueId: {
      type: 'string',
      description: 'Unique identifier',
    },
    enabledByDefault: {
      type: 'boolean',
      description: 'Whether entity is enabled by default',
    },
    qos: {
      type: 'integer',
      description: 'MQTT QoS level',
      minimum: 0,
      maximum: 2,
    },
    retain: {
      type: 'boolean',
      description: 'Retain MQTT messages',
    },
    device: {
      type: 'object',
      description: 'Device configuration',
      properties: {
        name: { type: 'string', description: 'Device name' },
        manufacturer: { type: 'string', description: 'Manufacturer' },
        model: { type: 'string', description: 'Model' },
        identifiers: {
          type: 'array',
          items: { type: 'string' },
          description: 'Device identifiers',
        },
      },
    },
    availability: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          topic: { type: 'string', description: 'Availability topic' },
          payloadAvailable: { type: 'string', description: 'Available payload' },
          payloadNotAvailable: { type: 'string', description: 'Not available payload' },
        },
      },
    },
  },
}

export const mockKubernetesButton = {
  apiVersion: 'mqtt.home-assistant.io/v1alpha1',
  kind: 'MQTTButton',
  metadata: {
    name: 'test-button',
    namespace: 'default',
    resourceVersion: '12345',
    creationTimestamp: '2024-01-15T10:00:00Z',
  },
  spec: {
    name: 'Test Button',
    commandTopic: 'homeassistant/button/test/command',
    payloadPress: 'PRESS',
    icon: 'mdi:button-pointer',
  },
  status: {
    lastPublished: '2024-01-15T10:01:00Z',
    discoveryTopic: 'homeassistant/button/default-test-button/config',
    conditions: [
      {
        type: 'Published',
        status: 'True',
        lastTransitionTime: '2024-01-15T10:01:00Z',
      },
    ],
  },
}
