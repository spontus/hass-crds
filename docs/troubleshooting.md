# Troubleshooting

Common issues and their solutions when using hass-crds.

## Controller Issues

### Controller Won't Start

**Symptoms**: Pod in `CrashLoopBackOff` or `Error` state.

**Check logs**:
```bash
kubectl logs -n hass-crds deployment/hass-crds-controller
```

**Common causes**:

| Error | Cause | Solution |
|-------|-------|----------|
| `MQTT_HOST is required` | Missing environment variable | Ensure the `mqtt-credentials` Secret exists and contains `MQTT_HOST` |
| `connection refused` | MQTT broker unreachable | Verify broker hostname/IP and port; check network policies |
| `authentication failed` | Wrong credentials | Verify username/password in the Secret |
| `certificate verify failed` | TLS certificate issue | See [TLS Troubleshooting](controller.md#troubleshooting-tls) |

### Controller Running But Not Publishing

**Symptoms**: Controller logs show no errors, but entities don't appear in Home Assistant.

**Checklist**:

1. **Verify MQTT connection**:
   ```bash
   kubectl logs -n hass-crds deployment/hass-crds-controller | grep -i "connected"
   ```
   Look for: `Connected to MQTT broker at <host>:<port>`

2. **Check CRD status**:
   ```bash
   kubectl get mqttbuttons -n hass-crds -o yaml
   ```
   Look for `.status.conditions` — `Published` should be `True`.

3. **Verify discovery prefix matches Home Assistant**:
   - Default is `homeassistant`
   - Check HA's MQTT integration settings: Settings > Devices & Services > MQTT > Configure

4. **Check if dry-run is enabled**:
   ```bash
   kubectl get mqttbutton <name> -o yaml | grep dry-run
   ```

## MQTT Issues

### Entities Not Appearing in Home Assistant

**Step 1**: Verify the discovery message is being published.

Use an MQTT client to subscribe to the discovery topic:
```bash
mosquitto_sub -h <broker> -t "homeassistant/#" -v
```

**Step 2**: Check if the message is retained.

Discovery messages must be retained. If using a fresh broker, ensure `retain=true` is working.

**Step 3**: Verify Home Assistant MQTT discovery is enabled.

In HA: Settings > Devices & Services > MQTT > Configure > Enable discovery.

### Entities Disappear After Broker Restart

**Cause**: Retained messages were lost (some brokers don't persist retained messages by default).

**Solution**:
- Configure your broker to persist retained messages
- Reduce `RECONCILE_INTERVAL` to re-publish more frequently
- Use `rediscoverInterval` on critical entities

### Duplicate Entities in Home Assistant

**Cause**: Multiple CRs with the same `uniqueId`, or residual discovery messages from deleted entities.

**Solution**:
1. Check for duplicate uniqueIds:
   ```bash
   kubectl get mqttbuttons,mqttswitches,mqttsensors -A -o jsonpath='{range .items[*]}{.spec.uniqueId}{"\n"}{end}' | sort | uniq -d
   ```

2. Clean up orphaned discovery messages:
   ```bash
   mosquitto_pub -h <broker> -t "homeassistant/button/<old-id>/config" -n -r
   ```

## CRD Issues

### Schema Validation Errors

**Symptoms**: `kubectl apply` fails with validation error.

**Example**:
```
error: error validating "button.yaml": error validating data:
ValidationError(MQTTButton.spec): missing required field "commandTopic"
```

**Solution**: Check the [CRD Reference](crds/README.md) for required fields.

### Status Not Updating

**Symptoms**: `.status.lastPublished` never updates.

**Causes**:
1. Controller doesn't have RBAC permissions for status subresources
2. Controller is watching a different namespace

**Check RBAC**:
```bash
kubectl auth can-i update mqttbuttons/status --as=system:serviceaccount:hass-crds:hass-crds-controller -n hass-crds
```

### Finalizer Stuck on Deletion

**Symptoms**: Resource stuck in `Terminating` state.

**Cause**: Controller can't reach MQTT broker to publish empty payload.

**Solution** (if controller is unavailable):
```bash
kubectl patch mqttbutton <name> -n hass-crds --type=json -p='[{"op": "remove", "path": "/metadata/finalizers"}]'
```

Note: This leaves the entity orphaned in Home Assistant.

## Home Assistant Issues

### Entity Shows "Unavailable"

**Causes**:
1. Availability topic not receiving expected payload
2. Device went offline

**Check availability configuration**:
```bash
kubectl get mqttsensor <name> -o yaml | grep -A5 availability
```

Verify the availability topic is receiving the expected payload:
```bash
mosquitto_sub -h <broker> -t "<availability-topic>" -v
```

### Entity State Not Updating

**Causes**:
1. State topic not configured
2. State topic not receiving messages
3. Value template error

**Debug with MQTT client**:
```bash
mosquitto_sub -h <broker> -t "<state-topic>" -v
```

### Entity Actions Not Working

**Causes**:
1. Command topic misconfigured
2. Device not subscribed to command topic
3. Payload format incorrect

**Test command manually**:
```bash
mosquitto_pub -h <broker> -t "<command-topic>" -m "ON"
```

## Getting Help

If you can't resolve the issue:

1. **Collect diagnostics**:
   ```bash
   kubectl logs -n hass-crds deployment/hass-crds-controller > controller.log
   kubectl get mqttbuttons,mqttswitches,mqttsensors -A -o yaml > crds.yaml
   ```

2. **Open an issue** at [GitHub Issues](https://github.com/spontus/hass-crds/issues) with:
   - Controller logs (redact sensitive data)
   - CRD manifests
   - Home Assistant version
   - MQTT broker type and version
