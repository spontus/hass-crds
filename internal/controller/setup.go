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

package controller

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spontus/hass-crds/internal/mqtt"
)

// SetupAllControllers sets up all MQTT entity controllers with the manager.
// Note: Currently only Button, Switch, and Sensor are implemented.
// Additional controllers follow the same pattern and can be added as needed.
func SetupAllControllers(mgr ctrl.Manager, mqttClient mqtt.Client, log logr.Logger) error {
	c := mgr.GetClient()
	scheme := mgr.GetScheme()

	// Setup MQTTButton controller
	if err := setupMQTTButtonController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}

	// Setup MQTTSwitch controller
	if err := setupMQTTSwitchController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}

	// Setup MQTTSensor controller
	if err := setupMQTTSensorController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}

	// TODO: Add remaining controllers as they are implemented
	// - MQTTBinarySensor
	// - MQTTNumber
	// - MQTTSelect
	// - MQTTText
	// - MQTTScene
	// - MQTTTag
	// - MQTTLight
	// - MQTTCover
	// - MQTTLock
	// - MQTTValve
	// - MQTTFan
	// - MQTTSiren
	// - MQTTCamera
	// - MQTTImage
	// - MQTTNotify
	// - MQTTUpdate
	// - MQTTClimate
	// - MQTTHumidifier
	// - MQTTWaterHeater
	// - MQTTVacuum
	// - MQTTLawnMower
	// - MQTTAlarmControlPanel
	// - MQTTDeviceTracker
	// - MQTTDeviceTrigger
	// - MQTTEvent

	return nil
}

func setupMQTTButtonController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	reconciler := NewMQTTButtonReconciler(c, scheme, log, mqttClient)
	return reconciler.SetupWithManager(mgr)
}

func setupMQTTSwitchController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	reconciler := NewMQTTSwitchReconciler(c, scheme, log, mqttClient)
	return reconciler.SetupWithManager(mgr)
}

func setupMQTTSensorController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	reconciler := NewMQTTSensorReconciler(c, scheme, log, mqttClient)
	return reconciler.SetupWithManager(mgr)
}
