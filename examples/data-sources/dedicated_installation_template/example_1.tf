data "ovh_dedicated_installation_template" "ovh_template" {
  template_name = "debian12_64"
}

output "template" {
  value = data.ovh_dedicated_installation_template.ovh_template
}
