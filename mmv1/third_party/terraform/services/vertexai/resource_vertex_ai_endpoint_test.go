package vertexai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVertexAIEndpoint_vertexAiEndpointNetwork(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"endpoint_name": fmt.Sprint(acctest.RandInt(t) % 9999999999),
		"kms_key_name":  acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vertex-ai-endpoint-update-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointNetwork(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointNetworkUpdate(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIEndpoint_vertexAiEndpointNetwork(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "%{endpoint_name}"
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
  predict_request_response_logging_config {
    bigquery_destination {
      output_uri = "bq://${data.google_project.project.project_id}.${google_bigquery_dataset.bq_dataset.dataset_id}.request_response_logging"
    }
    enabled       = true
    sampling_rate = 0.1
  }

  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "bq_dataset" {
  dataset_id                 = "some_dataset%{endpoint_name}"
  friendly_name              = "logging dataset"
  description                = "This is a dataset that requests are logged to"
  location                   = "US"
  delete_contents_on_destroy = true
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIEndpoint_vertexAiEndpointNetworkUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "%{endpoint_name}"
  display_name = "new-sample-endpoint"
  description  = "An updated sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-two = "value-two"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

data "google_project" "project" {}
`, context)
}

func TestAccVertexAIEndpoint_vertexAiEndpointPrivateServiceConnectExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointPrivateServiceConnect_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "location", "name", "region", "terraform_labels"},
			},
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointPrivateServiceConnect_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "location", "name", "region", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIEndpoint_vertexAiEndpointPrivateServiceConnect_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name = "network-%{random_suffix}"
}

resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "endpoint-name%{random_suffix}"
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  private_service_connect_config {
    enable_private_service_connect = true
    project_allowlist = [
      "${data.google_project.project.project_id}"
    ]

    psc_automation_configs {
      project_id = data.google_project.project.project_id
      network    = google_compute_network.default.id
    }
  }
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIEndpoint_vertexAiEndpointPrivateServiceConnect_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "basic" {
  project_id      = "tf-test-id%{random_suffix}"
  name            = "tf-test-id%{random_suffix}"
  org_id          = "%{org_id}"
  deletion_policy = "DELETE"
  billing_account = "%{billing_account}"
}

resource "time_sleep" "wait_2_mins" {
  create_duration = "120s"

  depends_on = [google_project.basic]
}

resource "google_project_service" "basic" {
  project = google_project.basic.project_id
  service = "compute.googleapis.com"
	
	depends_on = [time_sleep.wait_2_mins]
}

resource "google_compute_network" "default" {
	project = google_project.basic.project_id
  name    = "network-%{random_suffix}"
	
	depends_on = [google_project_service.basic]
}

resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "endpoint-name%{random_suffix}"
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  private_service_connect_config {
    enable_private_service_connect = true
    project_allowlist = [
      "${google_project.basic.project_id}"
    ]

    psc_automation_configs {
      project_id = google_project.basic.project_id
      network    = google_compute_network.default.id
    }
  }
}
`, context)
}
