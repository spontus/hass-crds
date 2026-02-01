# MQTTLawnMower

A robot lawn mower entity with start mowing, pause, and dock commands.

- **Kind**: `MQTTLawnMower`
- **Resource**: `mqttlawnmowers`
- **HA Component**: `lawn_mower`
- **HA Docs**: [MQTT Lawn Mower](https://www.home-assistant.io/integrations/lawn_mower.mqtt/)

## Type-Specific Fields

| CRD Field | MQTT Key | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `activityStateTopic` | `activity_state_topic` | `string` | No | -- | Topic to read mower activity state |
| `activityValueTemplate` | `activity_value_template` | `string` | No | -- | Template to extract activity from payload |
| `dockCommandTopic` | `dock_command_topic` | `string` | No | -- | Topic to publish dock command |
| `dockCommandTemplate` | `dock_command_template` | `string` | No | -- | Template for dock command payload |
| `pauseCommandTopic` | `pause_command_topic` | `string` | No | -- | Topic to publish pause command |
| `pauseCommandTemplate` | `pause_command_template` | `string` | No | -- | Template for pause command payload |
| `startMowingCommandTopic` | `start_mowing_command_topic` | `string` | No | -- | Topic to publish start mowing command |
| `startMowingCommandTemplate` | `start_mowing_command_template` | `string` | No | -- | Template for start mowing command payload |
| `optimistic` | `optimistic` | `bool` | No | `false` | Assume state changes immediately |

In addition to the fields above, all [common fields](common-fields.md) are supported.

## Example

```yaml
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLawnMower
metadata:
  name: robot-mower
spec:
  name: "Robot Mower"
  activityStateTopic: "mower/robot/activity"
  dockCommandTopic: "mower/robot/dock"
  pauseCommandTopic: "mower/robot/pause"
  startMowingCommandTopic: "mower/robot/start"
  device:
    name: "Robot Mower"
    identifiers:
      - "robot-mower-01"
    manufacturer: "Custom"
    model: "AutoMow v2"
    suggestedArea: "Garden"
```
