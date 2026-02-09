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

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch

func SetupAllControllers(mgr ctrl.Manager, mqttClient mqtt.Client, log logr.Logger) error {
	c := mgr.GetClient()
	scheme := mgr.GetScheme()

	if err := setupMQTTButtonController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTSwitchController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTSensorController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTBinarySensorController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTNumberController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTSelectController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTTextController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTSceneController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTTagController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTLightController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTCoverController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTLockController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTValveController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTFanController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTSirenController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTCameraController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTImageController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTNotifyController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTUpdateController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTClimateController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTHumidifierController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTWaterHeaterController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTVacuumController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTLawnMowerController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTAlarmControlPanelController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTDeviceTrackerController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTDeviceTriggerController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}
	if err := setupMQTTEventController(c, scheme, log, mqttClient, mgr); err != nil {
		return err
	}

	return nil
}

func setupMQTTButtonController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTButtonReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTSwitchController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTSwitchReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTSensorController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTSensorReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTBinarySensorController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTBinarySensorReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTNumberController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTNumberReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTSelectController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTSelectReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTTextController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTTextReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTSceneController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTSceneReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTTagController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTTagReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTLightController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTLightReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTCoverController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTCoverReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTLockController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTLockReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTValveController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTValveReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTFanController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTFanReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTSirenController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTSirenReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTCameraController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTCameraReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTImageController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTImageReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTNotifyController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTNotifyReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTUpdateController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTUpdateReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTClimateController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTClimateReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTHumidifierController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTHumidifierReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTWaterHeaterController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTWaterHeaterReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTVacuumController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTVacuumReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTLawnMowerController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTLawnMowerReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTAlarmControlPanelController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTAlarmControlPanelReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTDeviceTrackerController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTDeviceTrackerReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTDeviceTriggerController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTDeviceTriggerReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}

func setupMQTTEventController(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client, mgr ctrl.Manager) error {
	return NewMQTTEventReconciler(c, scheme, log, mqttClient).SetupWithManager(mgr)
}
