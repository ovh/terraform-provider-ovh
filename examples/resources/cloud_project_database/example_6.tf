# The advanced plan of MongoDB accepts 3 to 7 nodes: this cluster has 5.
# All nodes must be in the same region.
# Plans, flavors and number of nodes of MongoDB:
# https://docs.ovhcloud.com/en/guides/public-cloud/databases/mongodb-concept-capabilities
resource "ovh_cloud_project_database" "mongodb_advanced" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-advanced-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "advanced"
  flavor       = "db2-15"

  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
}
