package vpcaccess_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVPCAccessConnector_vpcAccessConnectorThroughput(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput(context),
			},
			{
				ResourceName:      "google_vpc_access_connector.connector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCAccessConnector_vpcAccessConnectorThroughput_combiningThroughputAndInstancesFields_permadiff(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput_bothThroughputAndInstances(context),
				// When both min_instance/max_instance and min_throughput/max_throughput fields are sent to the API
				// the API ignores the throughput field values. Instead the API returns values for min and max throughput
				// based on the value of min and max instances. The mismatch with the config causes a permadiff.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCAccessConnector_vpcAccessConnectorThroughput_usingThroughputOrInstancesLimits(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput_justThroughputFields(context),
				Check: resource.ComposeTestCheckFunc(
					// These fields are set by the config in this test step
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "min_throughput", "400"),
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "max_throughput", "800"),
					// These fields aren't set in the config; the API sets and returns values
					// based on the thoughput values provided
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "min_instances", "4"),
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "max_instances", "8"),
				),
			},
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput_justInstanceFields(context),
				Check: resource.ComposeTestCheckFunc(
					// These fields are set by the config in this test step
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "min_instances", "4"),
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "max_instances", "6"),
					// These fields aren't set in the config; the API sets and returns values
					// based on the instance limit values provided
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "min_throughput", "400"),
					resource.TestCheckResourceAttr(
						"google_vpc_access_connector.connector", "max_throughput", "600"),
				),
			},
		},
	})
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput_bothThroughputAndInstances(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  min_throughput = 400
  max_throughput = 1000
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput_justInstanceFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 5
  max_instances = 7
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput_justThroughputFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_throughput = 400
  max_throughput = 800
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}
