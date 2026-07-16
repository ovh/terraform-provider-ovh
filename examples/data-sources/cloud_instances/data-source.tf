data "ovh_cloud_instances" "instances" {
  service_name = "<Public cloud project id>"
}

# Names of the instances that are up and running in GRA11.
output "gra11_running_instance_names" {
  value = sort([
    for instance in data.ovh_cloud_instances.instances.instances :
    instance.name
    if instance.region == "GRA11" && instance.current_state.power_state == "ACTIVE"
  ])
}

# Instance id keyed by name, ready to be fed to another resource.
output "instance_ids_by_name" {
  value = {
    for instance in data.ovh_cloud_instances.instances.instances :
    instance.name => instance.id
  }
}
