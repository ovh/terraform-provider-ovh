resource "ovh_cloud_managed_analytics" "opensearch" {
  service_name            = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description             = "my-first-opensearch"
  engine                  = "opensearch"
  version                 = "2"
  plan                    = "business"
  opensearch_acls_enabled = true
  nodes {
    region  = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region  = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region  = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  flavor                  = "db1-30"
}
