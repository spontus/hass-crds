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
	"os/exec"
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
				cmd := exec.Command("kubectl", "get", "mqttbutton", "test-button",
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd = exec.Command("kubectl", "delete", "mqttbutton", "delete-test-button", "-n", "hass-crds-e2e")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
			cmd := exec.Command("kubectl", "apply", "-f", "-")
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
})
