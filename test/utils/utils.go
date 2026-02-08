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

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:golint,revive
)

const (
	KindClusterName    = "hass-crds-e2e"
	KindContext        = "kind-hass-crds-e2e"
	TestNamespace      = "hass-crds-e2e"
	ControllerImage    = "hass-crds-controller:e2e"
	KindConfigFile     = "test/e2e/kind-config.yaml"
	ManifestsDir       = "test/e2e/manifests"
	MQTTBrokerManifest = "test/e2e/manifests/mosquitto.yaml"
	HAManifest         = "test/e2e/manifests/homeassistant.yaml"
	ControllerManfest  = "test/e2e/manifests/controller.yaml"
	// Pre-generated JWT for e2e testing (matches auth storage in homeassistant.yaml)
	HAAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlMmVfdGVzdF90b2tlbl9pZF8xMjM0NSIsImlhdCI6MTcwNDA2NzIwMCwiZXhwIjoxODkzNDU2MDAwfQ.hx4RNm-QWqEO-Zl5D-0EF3xdrpqfnW7apUUyeVvvsHI"
)

// Kubectl returns an exec.Cmd for kubectl with the e2e context pre-configured
func Kubectl(args ...string) *exec.Cmd {
	fullArgs := append([]string{"--context", KindContext}, args...)
	return exec.Command("kubectl", fullArgs...)
}

// Run executes the provided command within the project directory
func Run(cmd *exec.Cmd) ([]byte, error) {
	dir, _ := GetProjectDir()
	cmd.Dir = dir

	if err := os.Chdir(cmd.Dir); err != nil {
		fmt.Fprintf(GinkgoWriter, "chdir dir: %s\n", err)
	}

	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	command := strings.Join(cmd.Args, " ")
	fmt.Fprintf(GinkgoWriter, "running: %s\n", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%s failed with error: (%v) %s", command, err, string(output))
	}

	return output, nil
}

// GetProjectDir returns the project root directory
func GetProjectDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return wd, err
	}
	// Handle both running from project root and test/e2e
	wd = strings.Replace(wd, "/test/e2e", "", -1)
	return wd, nil
}

// GetNonEmptyLines splits output by newlines and filters empty strings
func GetNonEmptyLines(output string) []string {
	var res []string
	elements := strings.Split(output, "\n")
	for _, element := range elements {
		if element != "" {
			res = append(res, element)
		}
	}
	return res
}

// CreateKindCluster creates a new Kind cluster for E2E tests
func CreateKindCluster() error {
	projectDir, _ := GetProjectDir()
	configPath := filepath.Join(projectDir, KindConfigFile)

	// Check if cluster already exists
	cmd := exec.Command("kind", "get", "clusters")
	output, _ := Run(cmd)
	if strings.Contains(string(output), KindClusterName) {
		fmt.Fprintf(GinkgoWriter, "Kind cluster %s already exists\n", KindClusterName)
		return nil
	}

	cmd = exec.Command("kind", "create", "cluster",
		"--name", KindClusterName,
		"--config", configPath,
		"--wait", "120s",
	)
	_, err := Run(cmd)
	return err
}

// DeleteKindCluster deletes the Kind cluster
func DeleteKindCluster() error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", KindClusterName)
	_, err := Run(cmd)
	return err
}

// BuildAndLoadControllerImage builds the controller image and loads it into Kind
func BuildAndLoadControllerImage() error {
	// Build the controller image
	cmd := exec.Command("docker", "build", "-t", ControllerImage, ".")
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("failed to build controller image: %w", err)
	}

	// Load image into Kind cluster
	cmd = exec.Command("kind", "load", "docker-image", ControllerImage, "--name", KindClusterName)
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("failed to load image into Kind: %w", err)
	}

	return nil
}

// InstallCRDs installs the CRDs into the cluster
func InstallCRDs() error {
	cmd := exec.Command("kubectl", "--context", KindContext, "apply", "-f", "config/crd/crds.yaml")
	_, err := Run(cmd)
	return err
}

// DeployMosquitto deploys the MQTT broker
func DeployMosquitto() error {
	projectDir, _ := GetProjectDir()
	manifestPath := filepath.Join(projectDir, MQTTBrokerManifest)

	cmd := exec.Command("kubectl", "--context", KindContext, "apply", "-f", manifestPath)
	if _, err := Run(cmd); err != nil {
		return err
	}

	// Wait for broker to be ready
	return WaitForDeployment("mosquitto", TestNamespace, 120*time.Second)
}

// DeployHomeAssistant deploys Home Assistant
func DeployHomeAssistant() error {
	projectDir, _ := GetProjectDir()
	manifestPath := filepath.Join(projectDir, HAManifest)

	cmd := exec.Command("kubectl", "--context", KindContext, "apply", "-f", manifestPath)
	if _, err := Run(cmd); err != nil {
		return err
	}

	// Wait for Home Assistant to be ready (takes longer to start)
	return WaitForDeployment("homeassistant", TestNamespace, 180*time.Second)
}

// DeployController deploys the hass-crds controller
func DeployController() error {
	projectDir, _ := GetProjectDir()
	manifestPath := filepath.Join(projectDir, ControllerManfest)

	cmd := exec.Command("kubectl", "--context", KindContext, "apply", "-f", manifestPath)
	if _, err := Run(cmd); err != nil {
		return err
	}

	// Wait for controller to be ready
	return WaitForDeployment("hass-crds-controller", TestNamespace, 60*time.Second)
}

// WaitForDeployment waits for a deployment to have available replicas
func WaitForDeployment(name, namespace string, timeout time.Duration) error {
	cmd := exec.Command("kubectl", "--context", KindContext, "rollout", "status",
		"deployment/"+name,
		"-n", namespace,
		"--timeout", fmt.Sprintf("%.0fs", timeout.Seconds()),
	)
	_, err := Run(cmd)
	return err
}

// WaitForMQTTConnection waits for HA to connect to MQTT
func WaitForMQTTConnection(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Check HA logs for MQTT connection
		cmd := exec.Command("kubectl", "--context", KindContext, "logs",
			"deployment/homeassistant",
			"-n", TestNamespace,
			"--tail", "100",
		)
		output, err := Run(cmd)
		if err == nil && strings.Contains(string(output), "MQTT") {
			fmt.Fprintf(GinkgoWriter, "Home Assistant MQTT connection established\n")
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for MQTT connection")
}

// CleanupTestResources deletes test CRDs from the namespace
func CleanupTestResources() error {
	// Delete all MQTT resources in test namespace
	resources := []string{
		"mqttbuttons",
		"mqttswitches",
		"mqttsensors",
		"mqttbinarysensors",
		"mqttnumbers",
		"mqttselects",
		"mqtttexts",
		"mqttscenes",
		"mqtttags",
		"mqttlights",
		"mqttcovers",
		"mqttlocks",
		"mqttvalves",
		"mqttfans",
		"mqttsirens",
		"mqttcameras",
		"mqttimages",
		"mqttnotifys",
		"mqttupdates",
		"mqttclimates",
		"mqtthumidifiers",
		"mqttwaterheaters",
		"mqttvacuums",
		"mqttlawnmowers",
		"mqttalarmcontrolpanels",
		"mqttdevicetrackers",
		"mqttdevicetriggers",
		"mqttevents",
	}

	for _, resource := range resources {
		cmd := exec.Command("kubectl", "--context", KindContext, "delete", resource, "--all",
			"-n", TestNamespace,
			"--ignore-not-found=true",
		)
		_, _ = Run(cmd)
	}

	return nil
}

// GetHAEntityState queries Home Assistant API for entity state
func GetHAEntityState(entityID string) (string, error) {
	cmd := exec.Command("kubectl", "--context", KindContext, "exec",
		"deployment/homeassistant",
		"-n", TestNamespace,
		"--",
		"curl", "-s",
		"-H", fmt.Sprintf("Authorization: Bearer %s", HAAccessToken),
		"-H", "Content-Type: application/json",
		fmt.Sprintf("http://localhost:8123/api/states/%s", entityID),
	)

	output, err := Run(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// WaitForHAEntity waits for an entity to appear in Home Assistant
func WaitForHAEntity(entityID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		state, err := GetHAEntityState(entityID)
		if err == nil && !strings.Contains(state, "not found") && len(state) > 2 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout waiting for entity %s", entityID)
}

// PublishMQTTMessage publishes a message to MQTT (for testing state topics)
func PublishMQTTMessage(topic, payload string) error {
	cmd := exec.Command("kubectl", "--context", KindContext, "exec",
		"deployment/mosquitto",
		"-n", TestNamespace,
		"--",
		"mosquitto_pub",
		"-h", "localhost",
		"-t", topic,
		"-m", payload,
	)
	_, err := Run(cmd)
	return err
}

// SubscribeMQTTMessage subscribes and captures one message from a topic
func SubscribeMQTTMessage(topic string, timeout time.Duration) (string, error) {
	cmd := exec.Command("kubectl", "--context", KindContext, "exec",
		"deployment/mosquitto",
		"-n", TestNamespace,
		"--",
		"mosquitto_sub",
		"-h", "localhost",
		"-t", topic,
		"-C", "1",
		"-W", fmt.Sprintf("%.0f", timeout.Seconds()),
	)
	output, err := Run(cmd)
	return string(output), err
}

// GetPodLogs retrieves logs from a deployment
func GetPodLogs(deployment, namespace string, tailLines int) (string, error) {
	cmd := exec.Command("kubectl", "--context", KindContext, "logs",
		"deployment/"+deployment,
		"-n", namespace,
		"--tail", fmt.Sprintf("%d", tailLines),
	)
	output, err := Run(cmd)
	return string(output), err
}

// HAEntityExists checks if an entity exists in Home Assistant
func HAEntityExists(entityID string) (bool, error) {
	state, err := GetHAEntityState(entityID)
	if err != nil {
		return false, err
	}
	// HA returns 404 as JSON: {"message": "Entity not found."}
	if strings.Contains(state, "not found") || strings.Contains(state, "404") {
		return false, nil
	}
	// Valid entity response contains "entity_id"
	return strings.Contains(state, "entity_id"), nil
}

// GetHAEntities lists all entities in Home Assistant
func GetHAEntities() (string, error) {
	cmd := exec.Command("kubectl", "--context", KindContext, "exec",
		"deployment/homeassistant",
		"-n", TestNamespace,
		"--",
		"curl", "-s",
		"-H", fmt.Sprintf("Authorization: Bearer %s", HAAccessToken),
		"-H", "Content-Type: application/json",
		"http://localhost:8123/api/states",
	)

	output, err := Run(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// WaitForHAEntityWithAttributes waits for an entity to appear with expected attributes
func WaitForHAEntityWithAttributes(entityID string, expectedAttrs map[string]string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		state, err := GetHAEntityState(entityID)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		if strings.Contains(state, "not found") {
			time.Sleep(2 * time.Second)
			continue
		}

		// Check for expected attributes in the response
		allFound := true
		for key, value := range expectedAttrs {
			if !strings.Contains(state, fmt.Sprintf(`"%s"`, key)) ||
				!strings.Contains(state, value) {
				allFound = false
				break
			}
		}

		if allFound {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout waiting for entity %s with attributes %v", entityID, expectedAttrs)
}
