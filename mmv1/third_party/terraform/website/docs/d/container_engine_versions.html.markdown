---
subcategory: "Kubernetes (Container) Engine"
description: |-
  Provides lists of available Google Kubernetes Engine versions for masters and nodes.
---

# google_container_engine_versions

Provides access to available Google Kubernetes Engine versions in a zone or region for a given project.

To get more information about GKE versions, see:
  * [The API reference](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations/getServerConfig)

-> If you are using the `google_container_engine_versions` datasource with a
regional cluster, ensure that you have provided a region as the `location` to
the datasource. A region can have a different set of supported versions than
its component zones, and not all zones in a region are guaranteed to
support the same version.

## Example Usage

```hcl
data "google_container_engine_versions" "central1b" {
  provider       = google-beta
  location       = "us-central1-b"
  version_prefix = "1.12."
}

resource "google_container_cluster" "foo" {
  name               = "terraform-test-cluster"
  location           = "us-central1-b"
  node_version       = data.google_container_engine_versions.central1b.latest_node_version
  initial_node_count = 1
}

output "stable_channel_default_version" {
  value = data.google_container_engine_versions.central1b.release_channel_default_version["STABLE"]
}

output "stable_channel_latest_version" {
  value = data.google_container_engine_versions.central1b.release_channel_latest_version["STABLE"]
}
```

## Argument Reference

The following arguments are supported:

* `location` (Optional) - The location (region or zone) to list versions for.
Must exactly match the location the cluster will be deployed in, or listed
versions may not be available. If `location`, `region`, and `zone` are not
specified, the provider-level zone must be set and is used instead.

* `project` (Optional) - ID of the project to list available cluster versions for. Should match the project the cluster will be deployed to.
  Defaults to the project that the provider is authenticated with.

* `version_prefix` (Optional) - If provided, Terraform will only return versions
that match the string prefix. For example, `1.11.` will match all `1.11` series
releases. Since this is just a string match, it's recommended that you append a
`.` after minor versions to ensure that prefixes such as `1.1` don't match
versions like `1.12.5-gke.10` accidentally. See [the docs on versioning schema](https://cloud.google.com/kubernetes-engine/versioning-and-upgrades#versioning_scheme)
for full details on how version strings are formatted.

## Attributes Reference

The following attributes are exported:

* `valid_master_versions` - A list of versions available in the given zone for use with master instances.
* `valid_node_versions` - A list of versions available in the given zone for use with node instances.
* `latest_master_version` - The latest version available in the given zone for use with master instances.
* `latest_node_version` - The latest version available in the given zone for use with node instances.
* `default_cluster_version` - Version of Kubernetes the service deploys by default.
* `release_channel_default_version` - A map from a release channel name to the channel's default version. See the docs on [available release channel names](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters#Cluster.Channel_1) for more details.
* `release_channel_latest_version` - A map from a release channel name to the channel's latest version. See the docs on [available release channel names](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters#Cluster.Channel_1) for more details.
* `release_channel_upgrade_target_version` - A map from a release channel name to the channel's auto upgrade target version. See the docs on [available release channel names](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters#Cluster.Channel_1) for more details.
