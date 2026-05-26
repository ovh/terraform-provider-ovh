data "ovh_vps_models" "models" {}

output "model_names" {
  value = data.ovh_vps_models.models.models[*].name
}
