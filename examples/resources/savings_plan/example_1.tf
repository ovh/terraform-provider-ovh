resource "ovh_savings_plan" "plan" {
  service_name = "<public cloud project ID>"
  flavor = "Rancher"
  period = "P1Y"
  size = 2
  display_name = "one_year_rancher_savings_plan"
  auto_renewal = true
}
