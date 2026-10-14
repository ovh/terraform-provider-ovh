# The number of nodes is the number of "nodes" blocks: the production plan of
# MongoDB requires exactly 3 nodes. All nodes must be in the same region.
# Plans, flavors and number of nodes for each engine:
# https://docs.ovhcloud.com/en/guides/public-cloud/databases/capabilities
resource "ovh_cloud_project_database" "mongodb_production" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-production-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "production"
  flavor       = "db2-2"

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
