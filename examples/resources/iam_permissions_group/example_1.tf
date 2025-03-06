# create a group allowing all actions in the category READ on VPSs
resource "ovh_iam_permissions_group" "read_vps" {
  name        = "read_vps"
  description = "Read access to vps"

  allow = [
    for act in data.ovh_iam_reference_actions.vps.actions : act.action if(contains(act.categories, "READ"))
  ]
}

data "ovh_iam_reference_actions" "vps" {
  type = "vps"
}
