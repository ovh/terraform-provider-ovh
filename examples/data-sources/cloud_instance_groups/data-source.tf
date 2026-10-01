data "ovh_cloud_instance_groups" "groups" {
  service_name = "<Public cloud project id>"
}

# Names of the groups whose members are spread across distinct hypervisors.
output "anti_affinity_group_names" {
  value = sort([
    for group in data.ovh_cloud_instance_groups.groups.instance_groups :
    group.name if group.policy == "ANTI_AFFINITY"
  ])
}

# Member count per group.
output "group_member_counts" {
  value = {
    for group in data.ovh_cloud_instance_groups.groups.instance_groups :
    group.name => length(group.current_state.members)
  }
}
