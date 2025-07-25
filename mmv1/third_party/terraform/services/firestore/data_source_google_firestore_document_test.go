package firestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDatasourceFirestoreDocument_simple(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceFirestoreDocument_simple(randomSuffix, orgId, "doc-id-1", "val1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_firestore_document.instance", "fields",
						"{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"val1\"}}}}}"),
					resource.TestCheckResourceAttr("data.google_firestore_document.instance",
						"id", fmt.Sprintf("projects/tf-test-%s/databases/(default)/documents/somenewcollection/doc-id-1", randomSuffix)),
					resource.TestCheckResourceAttr("data.google_firestore_document.instance",
						"name", fmt.Sprintf("projects/tf-test-%s/databases/(default)/documents/somenewcollection/doc-id-1", randomSuffix)),
					resource.TestCheckResourceAttr("data.google_firestore_document.instance",
						"collection", "somenewcollection"),
					resource.TestCheckResourceAttr("data.google_firestore_document.instance",
						"database", "(default)"),
					resource.TestCheckResourceAttrSet("data.google_firestore_document.instance", "path"),
					resource.TestCheckResourceAttrSet("data.google_firestore_document.instance", "create_time"),
					resource.TestCheckResourceAttrSet("data.google_firestore_document.instance", "update_time"),
				),
			},
		},
	})
}

func testAccDatasourceFirestoreDocument_simple_basicDeps(randomSuffix, orgId string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "tf-test-%s"
	name       = "tf-test-%s"
	org_id     = "%s"
	deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
	depends_on = [google_project.project]

	create_duration = "60s"
}

resource "google_project_service" "firestore" {
	project = google_project.project.project_id
	service = "firestore.googleapis.com"

	# Needed for CI tests for permissions to propagate, should not be needed for actual usage
	depends_on = [time_sleep.wait_60_seconds]
}

resource "google_firestore_database" "database" {
	project     = google_project.project.project_id
	name        = "(default)"
	location_id = "nam5"
	type        = "FIRESTORE_NATIVE"

	depends_on = [google_project_service.firestore]
}
`, randomSuffix, randomSuffix, orgId)
}

func testAccDatasourceFirestoreDocument_simple(randomSuffix, orgId, name, val string) string {
	return testAccDatasourceFirestoreDocument_simple_basicDeps(randomSuffix, orgId) + fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	project     = google_project.project.project_id
	database    = google_firestore_database.database.name
	collection  = "somenewcollection"
	document_id = "%s"
	fields      = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"%s\"}}}}}"
}

data "google_firestore_document" "instance" {
	project     = google_firestore_document.instance.project
	database    = google_firestore_document.instance.database
	collection  = google_firestore_document.instance.collection
	document_id = google_firestore_document.instance.document_id
}
`, name, val)
}
