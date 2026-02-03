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
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spontus/hass-crds/test/utils"
)

var skipClusterTeardown = false

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting hass-crds E2E test suite\n")
	RunSpecs(t, "hass-crds E2E Suite")
}

var _ = BeforeSuite(func() {
	// Check if we should skip cluster creation (for debugging)
	if os.Getenv("SKIP_CLUSTER_SETUP") == "true" {
		fmt.Fprintf(GinkgoWriter, "Skipping cluster setup (SKIP_CLUSTER_SETUP=true)\n")
		return
	}

	By("Creating Kind cluster")
	err := utils.CreateKindCluster()
	Expect(err).NotTo(HaveOccurred(), "Failed to create Kind cluster")

	By("Building and loading controller image")
	err = utils.BuildAndLoadControllerImage()
	Expect(err).NotTo(HaveOccurred(), "Failed to build/load controller image")

	By("Installing CRDs")
	err = utils.InstallCRDs()
	Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")

	By("Deploying Mosquitto MQTT broker")
	err = utils.DeployMosquitto()
	Expect(err).NotTo(HaveOccurred(), "Failed to deploy Mosquitto")

	By("Deploying Home Assistant")
	err = utils.DeployHomeAssistant()
	Expect(err).NotTo(HaveOccurred(), "Failed to deploy Home Assistant")

	By("Deploying hass-crds controller")
	err = utils.DeployController()
	Expect(err).NotTo(HaveOccurred(), "Failed to deploy controller")

	fmt.Fprintf(GinkgoWriter, "E2E test environment ready\n")
})

var _ = AfterSuite(func() {
	// Check if we should skip teardown (for debugging)
	if os.Getenv("SKIP_CLUSTER_TEARDOWN") == "true" || skipClusterTeardown {
		fmt.Fprintf(GinkgoWriter, "Skipping cluster teardown (SKIP_CLUSTER_TEARDOWN=true)\n")
		fmt.Fprintf(GinkgoWriter, "To access the cluster: kubectl cluster-info --context kind-%s\n", utils.KindClusterName)
		fmt.Fprintf(GinkgoWriter, "To delete manually: kind delete cluster --name %s\n", utils.KindClusterName)
		return
	}

	By("Deleting Kind cluster")
	err := utils.DeleteKindCluster()
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Warning: failed to delete Kind cluster: %v\n", err)
	}
})
