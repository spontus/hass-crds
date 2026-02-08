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

package e2e

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spontus/hass-crds/test/utils"
)

var _ = Describe("MQTT Discovery", func() {
	AfterEach(func() {
		_ = utils.CleanupTestResources()
	})

	Context("MQTTButton", func() {
		It("should publish discovery message when created", func() {
			By("Creating an MQTTButton resource")
			buttonYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: test-button
  namespace: hass-crds-e2e
spec:
  name: "E2E Test Button"
  commandTopic: "e2e/button/test/command"
  payloadPress: "PRESS"
  icon: "mdi:button-pointer"
  device:
    name: "E2E Test Device"
    identifiers:
      - "e2e-test-device-001"
    manufacturer: "hass-crds"
    model: "E2E Test"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(buttonYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for discovery message on MQTT")
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/button/hass-crds-e2e/test-button/config", 10*time.Second)
				if err != nil {
					fmt.Fprintf(GinkgoWriter, "MQTT subscribe error: %v\n", err)
					return false
				}

				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					fmt.Fprintf(GinkgoWriter, "JSON unmarshal error: %v\n", err)
					return false
				}

				Expect(payload["name"]).To(Equal("E2E Test Button"))
				Expect(payload["command_topic"]).To(Equal("e2e/button/test/command"))
				Expect(payload["payload_press"]).To(Equal("PRESS"))
				Expect(payload["unique_id"]).To(Equal("hass-crds-e2e-test-button"))

				device, ok := payload["device"].(map[string]interface{})
				Expect(ok).To(BeTrue(), "device should be a map")
				Expect(device["name"]).To(Equal("E2E Test Device"))
				Expect(device["manufacturer"]).To(Equal("hass-crds"))

				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())

			By("Verifying resource status was updated")
			Eventually(func() bool {
				cmd := utils.Kubectl("get", "mqttbutton", "test-button",
					"-n", "hass-crds-e2e", "-o", "jsonpath={.status.discoveryTopic}")
				out, err := utils.Run(cmd)
				if err != nil {
					return false
				}
				return strings.Contains(string(out), "homeassistant/button")
			}, 30*time.Second, 2*time.Second).Should(BeTrue())
		})

		It("should publish empty payload when deleted", func() {
			By("Creating an MQTTButton resource")
			buttonYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: delete-test-button
  namespace: hass-crds-e2e
spec:
  name: "Delete Test Button"
  commandTopic: "e2e/button/delete-test/command"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(buttonYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for initial discovery message")
			topic := "homeassistant/button/hass-crds-e2e/delete-test-button/config"
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage(topic, 10*time.Second)
				return err == nil && len(msg) > 0
			}, 60*time.Second, 5*time.Second).Should(BeTrue())

			By("Deleting the MQTTButton resource")
			cmd = utils.Kubectl("delete", "mqttbutton", "delete-test-button", "-n", "hass-crds-e2e")
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying retained message is cleared (entity removal)")
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage(topic, 5*time.Second)
				if err != nil {
					return true
				}
				return len(strings.TrimSpace(msg)) == 0
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTSensor", func() {
		It("should publish sensor discovery with all fields", func() {
			By("Creating an MQTTSensor resource")
			sensorYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: test-temperature
  namespace: hass-crds-e2e
spec:
  name: "E2E Temperature Sensor"
  stateTopic: "e2e/sensor/temperature/state"
  unitOfMeasurement: "°C"
  deviceClass: "temperature"
  stateClass: "measurement"
  valueTemplate: "{{ value_json.temperature }}"
  device:
    name: "E2E Sensor Hub"
    identifiers:
      - "e2e-sensor-hub-001"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(sensorYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying discovery message contents")
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/sensor/hass-crds-e2e/test-temperature/config", 10*time.Second)
				if err != nil {
					return false
				}

				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}

				Expect(payload["name"]).To(Equal("E2E Temperature Sensor"))
				Expect(payload["state_topic"]).To(Equal("e2e/sensor/temperature/state"))
				Expect(payload["unit_of_measurement"]).To(Equal("°C"))
				Expect(payload["device_class"]).To(Equal("temperature"))
				Expect(payload["state_class"]).To(Equal("measurement"))
				Expect(payload["value_template"]).To(Equal("{{ value_json.temperature }}"))

				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTSwitch", func() {
		It("should publish switch discovery with command and state topics", func() {
			By("Creating an MQTTSwitch resource")
			switchYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSwitch
metadata:
  name: test-switch
  namespace: hass-crds-e2e
spec:
  name: "E2E Test Switch"
  commandTopic: "e2e/switch/test/set"
  stateTopic: "e2e/switch/test/state"
  payloadOn: "ON"
  payloadOff: "OFF"
  stateOn: "ON"
  stateOff: "OFF"
  optimistic: false
  device:
    name: "E2E Switch Device"
    identifiers:
      - "e2e-switch-001"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(switchYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying discovery message contents")
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/switch/hass-crds-e2e/test-switch/config", 10*time.Second)
				if err != nil {
					return false
				}

				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}

				Expect(payload["name"]).To(Equal("E2E Test Switch"))
				Expect(payload["command_topic"]).To(Equal("e2e/switch/test/set"))
				Expect(payload["state_topic"]).To(Equal("e2e/switch/test/state"))
				Expect(payload["payload_on"]).To(Equal("ON"))
				Expect(payload["payload_off"]).To(Equal("OFF"))

				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTBinarySensor", func() {
		It("should publish binary sensor discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTBinarySensor
metadata:
  name: test-motion
  namespace: hass-crds-e2e
spec:
  name: "E2E Motion Sensor"
  stateTopic: "e2e/binary_sensor/motion/state"
  deviceClass: "motion"
  payloadOn: "detected"
  payloadOff: "clear"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/binary_sensor/hass-crds-e2e/test-motion/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Motion Sensor"))
				Expect(payload["device_class"]).To(Equal("motion"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTNumber", func() {
		It("should publish number discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTNumber
metadata:
  name: test-brightness
  namespace: hass-crds-e2e
spec:
  name: "E2E Brightness"
  commandTopic: "e2e/number/brightness/set"
  stateTopic: "e2e/number/brightness/state"
  min: 0
  max: 100
  step: 1
  mode: "slider"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/number/hass-crds-e2e/test-brightness/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Brightness"))
				Expect(payload["mode"]).To(Equal("slider"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTSelect", func() {
		It("should publish select discovery with options", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSelect
metadata:
  name: test-mode
  namespace: hass-crds-e2e
spec:
  name: "E2E Mode Select"
  commandTopic: "e2e/select/mode/set"
  stateTopic: "e2e/select/mode/state"
  options:
    - "auto"
    - "manual"
    - "eco"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/select/hass-crds-e2e/test-mode/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Mode Select"))
				options := payload["options"].([]interface{})
				Expect(options).To(HaveLen(3))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTText", func() {
		It("should publish text discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTText
metadata:
  name: test-text
  namespace: hass-crds-e2e
spec:
  name: "E2E Text Input"
  commandTopic: "e2e/text/input/set"
  stateTopic: "e2e/text/input/state"
  min: 1
  max: 100
  mode: "text"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/text/hass-crds-e2e/test-text/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Text Input"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTScene", func() {
		It("should publish scene discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTScene
metadata:
  name: test-scene
  namespace: hass-crds-e2e
spec:
  name: "E2E Movie Scene"
  commandTopic: "e2e/scene/movie/activate"
  payloadOn: "ACTIVATE"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/scene/hass-crds-e2e/test-scene/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Movie Scene"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTLight", func() {
		It("should publish light discovery with brightness", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLight
metadata:
  name: test-light
  namespace: hass-crds-e2e
spec:
  name: "E2E Light"
  commandTopic: "e2e/light/test/set"
  stateTopic: "e2e/light/test/state"
  brightnessCommandTopic: "e2e/light/test/brightness/set"
  brightnessStateTopic: "e2e/light/test/brightness/state"
  brightnessScale: 255
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/light/hass-crds-e2e/test-light/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Light"))
				Expect(payload["brightness_scale"]).To(BeEquivalentTo(255))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTCover", func() {
		It("should publish cover discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTCover
metadata:
  name: test-cover
  namespace: hass-crds-e2e
spec:
  name: "E2E Garage Door"
  commandTopic: "e2e/cover/garage/set"
  stateTopic: "e2e/cover/garage/state"
  deviceClass: "garage"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/cover/hass-crds-e2e/test-cover/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Garage Door"))
				Expect(payload["device_class"]).To(Equal("garage"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTLock", func() {
		It("should publish lock discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLock
metadata:
  name: test-lock
  namespace: hass-crds-e2e
spec:
  name: "E2E Front Door Lock"
  commandTopic: "e2e/lock/front/set"
  stateTopic: "e2e/lock/front/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/lock/hass-crds-e2e/test-lock/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Front Door Lock"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTFan", func() {
		It("should publish fan discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTFan
metadata:
  name: test-fan
  namespace: hass-crds-e2e
spec:
  name: "E2E Ceiling Fan"
  commandTopic: "e2e/fan/ceiling/set"
  stateTopic: "e2e/fan/ceiling/state"
  percentageCommandTopic: "e2e/fan/ceiling/speed/set"
  percentageStateTopic: "e2e/fan/ceiling/speed/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/fan/hass-crds-e2e/test-fan/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Ceiling Fan"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTClimate", func() {
		It("should publish climate discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTClimate
metadata:
  name: test-thermostat
  namespace: hass-crds-e2e
spec:
  name: "E2E Thermostat"
  temperatureCommandTopic: "e2e/climate/thermostat/temp/set"
  temperatureStateTopic: "e2e/climate/thermostat/temp/state"
  currentTemperatureTopic: "e2e/climate/thermostat/current"
  modeCommandTopic: "e2e/climate/thermostat/mode/set"
  modeStateTopic: "e2e/climate/thermostat/mode/state"
  modes:
    - "off"
    - "heat"
    - "cool"
    - "auto"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/climate/hass-crds-e2e/test-thermostat/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Thermostat"))
				modes := payload["modes"].([]interface{})
				Expect(modes).To(HaveLen(4))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTAlarmControlPanel", func() {
		It("should publish alarm panel discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTAlarmControlPanel
metadata:
  name: test-alarm
  namespace: hass-crds-e2e
spec:
  name: "E2E Alarm Panel"
  commandTopic: "e2e/alarm/panel/set"
  stateTopic: "e2e/alarm/panel/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/alarm_control_panel/hass-crds-e2e/test-alarm/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Alarm Panel"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTVacuum", func() {
		It("should publish vacuum discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTVacuum
metadata:
  name: test-vacuum
  namespace: hass-crds-e2e
spec:
  name: "E2E Robot Vacuum"
  commandTopic: "e2e/vacuum/robot/command"
  stateTopic: "e2e/vacuum/robot/state"
  supportedFeatures:
    - "start"
    - "stop"
    - "return_home"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/vacuum/hass-crds-e2e/test-vacuum/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Robot Vacuum"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTDeviceTracker", func() {
		It("should publish device tracker discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDeviceTracker
metadata:
  name: test-tracker
  namespace: hass-crds-e2e
spec:
  name: "E2E Phone Tracker"
  stateTopic: "e2e/device_tracker/phone/state"
  sourceType: "gps"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/device_tracker/hass-crds-e2e/test-tracker/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Phone Tracker"))
				Expect(payload["source_type"]).To(Equal("gps"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTEvent", func() {
		It("should publish event discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTEvent
metadata:
  name: test-event
  namespace: hass-crds-e2e
spec:
  name: "E2E Doorbell"
  stateTopic: "e2e/event/doorbell/state"
  eventTypes:
    - "press"
    - "double_press"
  deviceClass: "doorbell"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/event/hass-crds-e2e/test-event/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Doorbell"))
				Expect(payload["device_class"]).To(Equal("doorbell"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTUpdate", func() {
		It("should publish update discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTUpdate
metadata:
  name: test-update
  namespace: hass-crds-e2e
spec:
  name: "E2E Firmware Update"
  stateTopic: "e2e/update/firmware/state"
  commandTopic: "e2e/update/firmware/install"
  deviceClass: "firmware"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/update/hass-crds-e2e/test-update/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Firmware Update"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTSiren", func() {
		It("should publish siren discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSiren
metadata:
  name: test-siren
  namespace: hass-crds-e2e
spec:
  name: "E2E Alarm Siren"
  commandTopic: "e2e/siren/alarm/set"
  stateTopic: "e2e/siren/alarm/state"
  availableTones:
    - "fire"
    - "burglar"
    - "doorbell"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/siren/hass-crds-e2e/test-siren/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Alarm Siren"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTHumidifier", func() {
		It("should publish humidifier discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTHumidifier
metadata:
  name: test-humidifier
  namespace: hass-crds-e2e
spec:
  name: "E2E Humidifier"
  commandTopic: "e2e/humidifier/living/set"
  stateTopic: "e2e/humidifier/living/state"
  targetHumidityCommandTopic: "e2e/humidifier/living/humidity/set"
  targetHumidityStateTopic: "e2e/humidifier/living/humidity/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/humidifier/hass-crds-e2e/test-humidifier/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Humidifier"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTValve", func() {
		It("should publish valve discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTValve
metadata:
  name: test-valve
  namespace: hass-crds-e2e
spec:
  name: "E2E Water Valve"
  commandTopic: "e2e/valve/water/set"
  stateTopic: "e2e/valve/water/state"
  deviceClass: "water"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/valve/hass-crds-e2e/test-valve/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Water Valve"))
				Expect(payload["device_class"]).To(Equal("water"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTCamera", func() {
		It("should publish camera discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTCamera
metadata:
  name: test-camera
  namespace: hass-crds-e2e
spec:
  name: "E2E Camera"
  topic: "e2e/camera/front/image"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/camera/hass-crds-e2e/test-camera/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Camera"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("MQTTImage", func() {
		It("should publish image discovery", func() {
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTImage
metadata:
  name: test-image
  namespace: hass-crds-e2e
spec:
  name: "E2E Image"
  urlTopic: "e2e/image/snapshot/url"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage("homeassistant/image/hass-crds-e2e/test-image/config", 10*time.Second)
				if err != nil {
					return false
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}
				Expect(payload["name"]).To(Equal("E2E Image"))
				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("Availability", func() {
		It("should include availability configuration in discovery", func() {
			By("Creating an MQTTSensor with availability")
			sensorYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: sensor-with-availability
  namespace: hass-crds-e2e
spec:
  name: "Sensor With Availability"
  stateTopic: "e2e/sensor/avail/state"
  availability:
    - topic: "e2e/sensor/avail/status"
      payloadAvailable: "online"
      payloadNotAvailable: "offline"
  availabilityMode: "all"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(sensorYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying availability in discovery message")
			availTopic := "homeassistant/sensor/hass-crds-e2e/sensor-with-availability/config"
			Eventually(func() bool {
				msg, err := utils.SubscribeMQTTMessage(availTopic, 10*time.Second)
				if err != nil {
					return false
				}

				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &payload); err != nil {
					return false
				}

				avail, ok := payload["availability"].([]interface{})
				if !ok || len(avail) == 0 {
					return false
				}

				availConfig, ok := avail[0].(map[string]interface{})
				if !ok {
					return false
				}

				Expect(availConfig["topic"]).To(Equal("e2e/sensor/avail/status"))
				Expect(availConfig["payload_available"]).To(Equal("online"))
				Expect(availConfig["payload_not_available"]).To(Equal("offline"))
				Expect(payload["availability_mode"]).To(Equal("all"))

				return true
			}, 60*time.Second, 5*time.Second).Should(BeTrue())
		})
	})

	Context("Controller Logs", func() {
		It("should log discovery publish events", func() {
			By("Creating a test resource")
			buttonYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: log-test-button
  namespace: hass-crds-e2e
spec:
  name: "Log Test Button"
  commandTopic: "e2e/button/log-test/command"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(buttonYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Checking controller logs for publish event")
			Eventually(func() bool {
				logs, err := utils.GetPodLogs("hass-crds-controller", "hass-crds-e2e", 50)
				if err != nil {
					return false
				}
				return strings.Contains(logs, "Published discovery message") ||
					strings.Contains(logs, "log-test-button")
			}, 30*time.Second, 2*time.Second).Should(BeTrue())
		})
	})

	Context("Home Assistant Entity Creation", func() {
		It("should create a button entity in Home Assistant", func() {
			By("Creating an MQTTButton resource")
			buttonYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTButton
metadata:
  name: ha-verify-button
  namespace: hass-crds-e2e
spec:
  name: "HA Verify Button"
  commandTopic: "e2e/button/ha-verify/command"
  icon: "mdi:button-pointer"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(buttonYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for entity to appear in Home Assistant")
			// Entity ID is derived from the name field: "HA Verify Button" → "ha_verify_button"
			entityID := "button.ha_verify_button"
			Eventually(func() bool {
				exists, err := utils.HAEntityExists(entityID)
				if err != nil {
					fmt.Fprintf(GinkgoWriter, "HA entity check error: %v\n", err)
					return false
				}
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue(), "Entity should exist in Home Assistant")

			By("Verifying entity attributes")
			state, err := utils.GetHAEntityState(entityID)
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "HA Entity State: %s\n", state)

			Expect(state).To(ContainSubstring(`"friendly_name"`))
			Expect(state).To(ContainSubstring("HA Verify Button"))
		})

		It("should create a sensor entity in Home Assistant", func() {
			By("Creating an MQTTSensor resource")
			sensorYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSensor
metadata:
  name: ha-verify-sensor
  namespace: hass-crds-e2e
spec:
  name: "HA Verify Temperature"
  stateTopic: "e2e/sensor/ha-verify/state"
  unitOfMeasurement: "°C"
  deviceClass: "temperature"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(sensorYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for sensor entity in Home Assistant")
			// Entity ID is derived from the name field: "HA Verify Temperature" → "ha_verify_temperature"
			entityID := "sensor.ha_verify_temperature"
			Eventually(func() bool {
				exists, err := utils.HAEntityExists(entityID)
				if err != nil {
					fmt.Fprintf(GinkgoWriter, "HA entity check error: %v\n", err)
					return false
				}
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue(), "Sensor entity should exist in Home Assistant")

			By("Publishing a state value and verifying HA receives it")
			err = utils.PublishMQTTMessage("e2e/sensor/ha-verify/state", "23.5")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				state, err := utils.GetHAEntityState(entityID)
				if err != nil {
					return false
				}
				return strings.Contains(state, "23.5")
			}, 30*time.Second, 2*time.Second).Should(BeTrue(), "Sensor should show published value")
		})

		It("should remove entity from Home Assistant when CRD is deleted", func() {
			By("Creating an MQTTSwitch resource")
			switchYAML := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSwitch
metadata:
  name: ha-delete-test
  namespace: hass-crds-e2e
spec:
  name: "HA Delete Test Switch"
  commandTopic: "e2e/switch/ha-delete/set"
  stateTopic: "e2e/switch/ha-delete/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(switchYAML)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for entity to appear in Home Assistant")
			// Entity ID is derived from the name field: "HA Delete Test Switch" → "ha_delete_test_switch"
			entityID := "switch.ha_delete_test_switch"
			Eventually(func() bool {
				exists, err := utils.HAEntityExists(entityID)
				return err == nil && exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Deleting the CRD")
			cmd = utils.Kubectl("delete", "mqttswitch", "ha-delete-test", "-n", "hass-crds-e2e")
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying entity is removed from Home Assistant")
			Eventually(func() bool {
				exists, err := utils.HAEntityExists(entityID)
				if err != nil {
					return false
				}
				return !exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue(), "Entity should be removed from Home Assistant")
		})
	})

	Context("Home Assistant State Verification", func() {
		It("should update switch state when MQTT state is published", func() {
			By("Creating an MQTTSwitch resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSwitch
metadata:
  name: state-test-switch
  namespace: hass-crds-e2e
spec:
  name: "State Test Switch"
  commandTopic: "e2e/switch/state-test/set"
  stateTopic: "e2e/switch/state-test/state"
  payloadOn: "ON"
  payloadOff: "OFF"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "switch.state_test_switch"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing ON state and verifying")
			err = utils.PublishMQTTMessage("e2e/switch/state-test/state", "ON")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing OFF state and verifying")
			err = utils.PublishMQTTMessage("e2e/switch/state-test/state", "OFF")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"off"`))
		})

		It("should update binary sensor state when MQTT state is published", func() {
			By("Creating an MQTTBinarySensor resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTBinarySensor
metadata:
  name: state-test-motion
  namespace: hass-crds-e2e
spec:
  name: "State Test Motion"
  stateTopic: "e2e/binary_sensor/state-test/state"
  deviceClass: "motion"
  payloadOn: "detected"
  payloadOff: "clear"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "binary_sensor.state_test_motion"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing detected state and verifying")
			err = utils.PublishMQTTMessage("e2e/binary_sensor/state-test/state", "detected")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing clear state and verifying")
			err = utils.PublishMQTTMessage("e2e/binary_sensor/state-test/state", "clear")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"off"`))
		})

		It("should update number state when MQTT state is published", func() {
			By("Creating an MQTTNumber resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTNumber
metadata:
  name: state-test-brightness
  namespace: hass-crds-e2e
spec:
  name: "State Test Brightness"
  commandTopic: "e2e/number/state-test/set"
  stateTopic: "e2e/number/state-test/state"
  min: 0
  max: 100
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "number.state_test_brightness"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing numeric state and verifying")
			err = utils.PublishMQTTMessage("e2e/number/state-test/state", "75")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"75"`))

			By("Publishing different value and verifying")
			err = utils.PublishMQTTMessage("e2e/number/state-test/state", "25")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"25"`))
		})

		It("should update select state when MQTT state is published", func() {
			By("Creating an MQTTSelect resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSelect
metadata:
  name: state-test-mode
  namespace: hass-crds-e2e
spec:
  name: "State Test Mode"
  commandTopic: "e2e/select/state-test/set"
  stateTopic: "e2e/select/state-test/state"
  options:
    - "auto"
    - "manual"
    - "eco"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "select.state_test_mode"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing 'eco' state and verifying")
			err = utils.PublishMQTTMessage("e2e/select/state-test/state", "eco")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"eco"`))

			By("Publishing 'manual' state and verifying")
			err = utils.PublishMQTTMessage("e2e/select/state-test/state", "manual")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"manual"`))
		})

		It("should update text state when MQTT state is published", func() {
			By("Creating an MQTTText resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTText
metadata:
  name: state-test-text
  namespace: hass-crds-e2e
spec:
  name: "State Test Text"
  commandTopic: "e2e/text/state-test/set"
  stateTopic: "e2e/text/state-test/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "text.state_test_text"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing text state and verifying")
			err = utils.PublishMQTTMessage("e2e/text/state-test/state", "Hello World")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"Hello World"`))
		})

		It("should update light state when MQTT state is published", func() {
			By("Creating an MQTTLight resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLight
metadata:
  name: state-test-light
  namespace: hass-crds-e2e
spec:
  name: "State Test Light"
  commandTopic: "e2e/light/state-test/set"
  stateTopic: "e2e/light/state-test/state"
  payloadOn: "ON"
  payloadOff: "OFF"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "light.state_test_light"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing ON state and verifying")
			err = utils.PublishMQTTMessage("e2e/light/state-test/state", "ON")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing OFF state and verifying")
			err = utils.PublishMQTTMessage("e2e/light/state-test/state", "OFF")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"off"`))
		})

		It("should update cover state when MQTT state is published", func() {
			By("Creating an MQTTCover resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTCover
metadata:
  name: state-test-cover
  namespace: hass-crds-e2e
spec:
  name: "State Test Cover"
  commandTopic: "e2e/cover/state-test/set"
  stateTopic: "e2e/cover/state-test/state"
  deviceClass: "garage"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "cover.state_test_cover"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing open state and verifying")
			err = utils.PublishMQTTMessage("e2e/cover/state-test/state", "open")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"open"`))

			By("Publishing closed state and verifying")
			err = utils.PublishMQTTMessage("e2e/cover/state-test/state", "closed")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"closed"`))
		})

		It("should update lock state when MQTT state is published", func() {
			By("Creating an MQTTLock resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTLock
metadata:
  name: state-test-lock
  namespace: hass-crds-e2e
spec:
  name: "State Test Lock"
  commandTopic: "e2e/lock/state-test/set"
  stateTopic: "e2e/lock/state-test/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "lock.state_test_lock"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing LOCKED state and verifying")
			err = utils.PublishMQTTMessage("e2e/lock/state-test/state", "LOCKED")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"locked"`))

			By("Publishing UNLOCKED state and verifying")
			err = utils.PublishMQTTMessage("e2e/lock/state-test/state", "UNLOCKED")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"unlocked"`))
		})

		It("should update fan state when MQTT state is published", func() {
			By("Creating an MQTTFan resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTFan
metadata:
  name: state-test-fan
  namespace: hass-crds-e2e
spec:
  name: "State Test Fan"
  commandTopic: "e2e/fan/state-test/set"
  stateTopic: "e2e/fan/state-test/state"
  percentageCommandTopic: "e2e/fan/state-test/percentage/set"
  percentageStateTopic: "e2e/fan/state-test/percentage/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "fan.state_test_fan"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing ON state and verifying")
			err = utils.PublishMQTTMessage("e2e/fan/state-test/state", "ON")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing percentage and verifying")
			err = utils.PublishMQTTMessage("e2e/fan/state-test/percentage/state", "75")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"percentage":75`))
		})

		It("should update valve state when MQTT state is published", func() {
			By("Creating an MQTTValve resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTValve
metadata:
  name: state-test-valve
  namespace: hass-crds-e2e
spec:
  name: "State Test Valve"
  commandTopic: "e2e/valve/state-test/set"
  stateTopic: "e2e/valve/state-test/state"
  deviceClass: "water"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "valve.state_test_valve"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing open state and verifying")
			err = utils.PublishMQTTMessage("e2e/valve/state-test/state", "open")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"open"`))

			By("Publishing closed state and verifying")
			err = utils.PublishMQTTMessage("e2e/valve/state-test/state", "closed")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"closed"`))
		})

		It("should update siren state when MQTT state is published", func() {
			By("Creating an MQTTSiren resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTSiren
metadata:
  name: state-test-siren
  namespace: hass-crds-e2e
spec:
  name: "State Test Siren"
  commandTopic: "e2e/siren/state-test/set"
  stateTopic: "e2e/siren/state-test/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "siren.state_test_siren"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing ON state and verifying")
			err = utils.PublishMQTTMessage("e2e/siren/state-test/state", "ON")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing OFF state and verifying")
			err = utils.PublishMQTTMessage("e2e/siren/state-test/state", "OFF")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"off"`))
		})

		It("should update alarm control panel state when MQTT state is published", func() {
			By("Creating an MQTTAlarmControlPanel resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTAlarmControlPanel
metadata:
  name: state-test-alarm
  namespace: hass-crds-e2e
spec:
  name: "State Test Alarm"
  commandTopic: "e2e/alarm/state-test/set"
  stateTopic: "e2e/alarm/state-test/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "alarm_control_panel.state_test_alarm"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing armed_home state and verifying")
			err = utils.PublishMQTTMessage("e2e/alarm/state-test/state", "armed_home")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"armed_home"`))

			By("Publishing disarmed state and verifying")
			err = utils.PublishMQTTMessage("e2e/alarm/state-test/state", "disarmed")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"disarmed"`))
		})

		It("should update device tracker state when MQTT state is published", func() {
			By("Creating an MQTTDeviceTracker resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTDeviceTracker
metadata:
  name: state-test-tracker
  namespace: hass-crds-e2e
spec:
  name: "State Test Tracker"
  stateTopic: "e2e/device_tracker/state-test/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "device_tracker.state_test_tracker"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing home state and verifying")
			err = utils.PublishMQTTMessage("e2e/device_tracker/state-test/state", "home")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"home"`))

			By("Publishing not_home state and verifying")
			err = utils.PublishMQTTMessage("e2e/device_tracker/state-test/state", "not_home")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"not_home"`))
		})

		It("should update image entity when URL is published", func() {
			By("Creating an MQTTImage resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTImage
metadata:
  name: state-test-image
  namespace: hass-crds-e2e
spec:
  name: "State Test Image"
  urlTopic: "e2e/image/state-test/url"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "image.state_test_image"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing image URL and verifying entity updates")
			err = utils.PublishMQTTMessage("e2e/image/state-test/url", "http://example.com/image.jpg")
			Expect(err).NotTo(HaveOccurred())

			// Image entity state changes from "unknown" when URL is received
			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).ShouldNot(ContainSubstring(`"state":"unknown"`))
		})

		It("should update humidifier state when MQTT state is published", func() {
			By("Creating an MQTTHumidifier resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTHumidifier
metadata:
  name: state-test-humidifier
  namespace: hass-crds-e2e
spec:
  name: "State Test Humidifier"
  commandTopic: "e2e/humidifier/state-test/set"
  stateTopic: "e2e/humidifier/state-test/state"
  targetHumidityCommandTopic: "e2e/humidifier/state-test/humidity/set"
  targetHumidityStateTopic: "e2e/humidifier/state-test/humidity/state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "humidifier.state_test_humidifier"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing ON state and verifying")
			err = utils.PublishMQTTMessage("e2e/humidifier/state-test/state", "ON")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"on"`))

			By("Publishing target humidity and verifying")
			err = utils.PublishMQTTMessage("e2e/humidifier/state-test/humidity/state", "65")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"humidity":65`))
		})

		It("should update climate state when MQTT state is published", func() {
			By("Creating an MQTTClimate resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTClimate
metadata:
  name: state-test-climate
  namespace: hass-crds-e2e
spec:
  name: "State Test Climate"
  modeCommandTopic: "e2e/climate/state-test/mode/set"
  modeStateTopic: "e2e/climate/state-test/mode/state"
  temperatureCommandTopic: "e2e/climate/state-test/temp/set"
  temperatureStateTopic: "e2e/climate/state-test/temp/state"
  currentTemperatureTopic: "e2e/climate/state-test/current"
  modes:
    - "off"
    - "heat"
    - "cool"
    - "auto"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "climate.state_test_climate"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing heat mode and verifying")
			err = utils.PublishMQTTMessage("e2e/climate/state-test/mode/state", "heat")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"heat"`))

			By("Publishing current temperature and verifying")
			err = utils.PublishMQTTMessage("e2e/climate/state-test/current", "21.5")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"current_temperature":21.5`))

			By("Publishing target temperature and verifying")
			err = utils.PublishMQTTMessage("e2e/climate/state-test/temp/state", "23")
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"temperature":23`))
		})

		It("should update vacuum state when MQTT state is published", func() {
			By("Creating an MQTTVacuum resource")
			yaml := `
apiVersion: mqtt.home-assistant.io/v1alpha1
kind: MQTTVacuum
metadata:
  name: state-test-vacuum
  namespace: hass-crds-e2e
spec:
  name: "State Test Vacuum"
  commandTopic: "e2e/vacuum/state-test/command"
  stateTopic: "e2e/vacuum/state-test/state"
  schema: "state"
`
			cmd := utils.Kubectl("apply", "-f", "-")
			cmd.Stdin = strings.NewReader(yaml)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			entityID := "vacuum.state_test_vacuum"
			Eventually(func() bool {
				exists, _ := utils.HAEntityExists(entityID)
				return exists
			}, 90*time.Second, 5*time.Second).Should(BeTrue())

			By("Publishing cleaning state via JSON and verifying")
			err = utils.PublishMQTTMessage("e2e/vacuum/state-test/state", `{"state": "cleaning"}`)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"cleaning"`))

			By("Publishing docked state via JSON and verifying")
			err = utils.PublishMQTTMessage("e2e/vacuum/state-test/state", `{"state": "docked"}`)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				state, _ := utils.GetHAEntityState(entityID)
				return state
			}, 30*time.Second, 2*time.Second).Should(ContainSubstring(`"state":"docked"`))
		})
	})
})
