resource "ovh_cloud_project_rancher" "rancher" {
  project_id = "<public cloud project ID>"
  target_spec = {
    name = "MyRancher"
    plan = "STANDARD"
  }
}

output "rancher_url" {
  value = ovh_cloud_project_rancher.rancher.current_state.url
}
