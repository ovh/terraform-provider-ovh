data "ovh_cloud_managed_analytics" "m3db" {
  service_name  = "XXX"
  engine        = "m3db"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_analytics_m3db_namespace" "namespace" {
  service_name              = data.ovh_cloud_managed_analytics.m3db.service_name
  cluster_id                = data.ovh_cloud_managed_analytics.m3db.id
  name                      = "mynamespace"
  resolution                = "P2D"
  retention_period_duration = "PT48H"
}
