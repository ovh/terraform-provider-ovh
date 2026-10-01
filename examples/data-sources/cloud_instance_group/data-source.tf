data "ovh_cloud_instance_group" "group" {
  service_name = "<Public cloud project id>"
  id           = "<instance group id>"
}

# Membership is fixed at instance-creation time through the instance's
# group_id; this is the only way to read it back.
output "group_member_ids" {
  value = [
    for member in data.ovh_cloud_instance_group.group.current_state.members :
    member.id
  ]
}
