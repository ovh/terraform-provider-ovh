resource "ovh_vps_abort_snapshot" "abort" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Change any value here to re-run the abort.
  triggers = {
    run = "2026-05-25T00:00:00Z"
  }
}
